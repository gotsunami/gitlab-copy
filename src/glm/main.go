package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/xanzy/go-gitlab"
)

const (
	resultsPerPage = 100
)

// GitLab server endpoints
type endpoint struct {
	from, to *gitlab.Client
}

type migration struct {
	params   *config
	endpoint *endpoint
}

func NewMigration(c *config) (*migration, error) {
	if c == nil {
		return nil, errors.New("nil params")
	}
	fromgl := gitlab.NewClient(nil, c.From.Token)
	if err := fromgl.SetBaseURL(c.From.ServerURL); err != nil {
		return nil, err
	}
	togl := gitlab.NewClient(nil, c.To.Token)
	if err := togl.SetBaseURL(c.To.ServerURL); err != nil {
		return nil, err
	}
	m := &migration{c, &endpoint{fromgl, togl}}
	return m, nil
}

// Returns project by name.
func (m *migration) project(endpoint *gitlab.Client, name string) (*gitlab.Project, error) {
	proj, resp, err := endpoint.Projects.GetProject(name)
	if resp == nil {
		return nil, errors.New("network error: " + err.Error())
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("source project '%s' not found", name)
	}
	if err != nil {
		return nil, err
	}
	return proj, nil
}

func (m *migration) sourceProject(name string) (*gitlab.Project, error) {
	return m.project(m.endpoint.from, name)
}

func (m *migration) destProject(name string) (*gitlab.Project, error) {
	return m.project(m.endpoint.to, name)
}

func map2human(m map[string]int) string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return strings.Join(keys, ", ")
}

type projectStats struct {
	project                               *gitlab.Project
	nbIssues, nbClosed, nbOpened, nbNotes int
	milestones, labels                    map[string]int
}

func newProjectStats(prj *gitlab.Project) *projectStats {
	p := new(projectStats)
	p.project = prj
	p.milestones = make(map[string]int)
	p.labels = make(map[string]int)
	return p
}

func (p *projectStats) String() string {
	return fmt.Sprintf("%d issues (%d opened, %d closed)", p.nbIssues, p.nbOpened, p.nbClosed)
}

func (p *projectStats) pagination(client *gitlab.Client, f func(*gitlab.Client, *gitlab.ListOptions) (bool, error)) error {
	if client == nil {
		return errors.New("nil client")
	}

	curPage := 1
	opts := &gitlab.ListOptions{PerPage: resultsPerPage, Page: curPage}

	for {
		stop, err := f(client, opts)
		if err != nil {
			return err
		}
		if stop {
			break
		}
		curPage++
		opts.Page = curPage
	}
	return nil
}

func (p *projectStats) computeStats(client *gitlab.Client) error {
	if client == nil {
		return errors.New("nil client")
	}

	action := func(c *gitlab.Client, lo *gitlab.ListOptions) (bool, error) {
		opts := &gitlab.ListProjectIssuesOptions{ListOptions: gitlab.ListOptions{PerPage: lo.PerPage, Page: lo.Page}}
		issues, _, err := client.Issues.ListProjectIssues(*p.project.ID, opts)
		if err != nil {
			return false, err
		}
		if len(issues) > 0 {
			p.nbIssues += len(issues)
			for _, issue := range issues {
				switch issue.State {
				case "opened":
					p.nbOpened++
				case "closed":
					p.nbClosed++
				}
				if issue.Milestone.Title != "" {
					p.milestones[issue.Milestone.Title]++
				}
				if len(issue.Labels) > 0 {
					for _, label := range issue.Labels {
						p.labels[label]++
					}
				}
			}
		} else {
			// Exit
			return true, nil
		}
		return false, nil
	}

	if err := p.pagination(client, action); err != nil {
		return err
	}
	return nil
}

func (p *projectStats) computeIssueNotes(client *gitlab.Client) error {
	if client == nil {
		return errors.New("nil client")
	}

	action := func(c *gitlab.Client, lo *gitlab.ListOptions) (bool, error) {
		opts := &gitlab.ListProjectIssuesOptions{ListOptions: gitlab.ListOptions{PerPage: lo.PerPage, Page: lo.Page}}
		issues, _, err := client.Issues.ListProjectIssues(*p.project.ID, opts)
		if err != nil {
			return false, err
		}
		if len(issues) > 0 {
			for _, issue := range issues {
				notes, _, err := client.Notes.ListIssueNotes(*p.project.ID, issue.ID, nil)
				if err != nil {
					return false, err
				}
				p.nbNotes += len(notes)
			}
		} else {
			// Exit
			return true, nil
		}
		return false, nil
	}

	if err := p.pagination(client, action); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("missing config file")
	}
	c, err := parseConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	m, err := NewMigration(c)
	if err != nil {
		log.Fatal(err)
	}
	srcproj, err := m.sourceProject(c.From.Name)
	if err != nil {
		log.Fatal(err)
	}
	if srcproj == nil {
		log.Fatalf("source project not found on %s", c.From.ServerURL)
	}
	fmt.Printf("source: %s at %s\n", c.From.Name, c.From.ServerURL)

	dstproj, err := m.destProject(c.To.Name)
	if err != nil {
		log.Fatal(err)
	}
	if dstproj == nil {
		log.Fatalf("target project not found on %s", c.To.ServerURL)
	}
	fmt.Printf("target: %s at %s\n", c.To.Name, c.To.ServerURL)
	fmt.Println("--")

	// Find out how many issues we have
	fmt.Printf("source: finding issues ... ")

	pstats := newProjectStats(srcproj)

	if err := pstats.computeStats(m.endpoint.from); err != nil {
		log.Fatal(err)
	}
	fmt.Println("OK")
	fmt.Printf("source: %v\n", pstats)
	if len(pstats.milestones) > 0 {
		fmt.Printf("source: %d milestone(s): %s\n", len(pstats.milestones), map2human(pstats.milestones))
	}
	if len(pstats.labels) > 0 {
		fmt.Printf("source: %d label(s): %s\n", len(pstats.labels), map2human(pstats.labels))
	}

	fmt.Printf("source: counting notes (comments), can take a while ... ")
	if err := pstats.computeIssueNotes(m.endpoint.from); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\rsource: %d notes%50s\n", pstats.nbNotes, " ")
	fmt.Println("--")
	fmt.Println(`Migration rules are:
- Create milestone if not existing on target
- Create label if not existing on target
- Create issue if not existing on target (by title), either closed of opened on source
- Creaate note (attached to issue) if not existing on target

Use the --apply option parameter if that looks good to you to start the issues migration.
`)

}
