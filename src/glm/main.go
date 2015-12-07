package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/xanzy/go-gitlab"
)

const (
	resultsPerPage = 50
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

	curPage := 1
	nbIssues := 0
	nbClosed, nbOpened := 0, 0
	milestones := make(map[string]int)
	opts := &gitlab.ListProjectIssuesOptions{ListOptions: gitlab.ListOptions{PerPage: resultsPerPage, Page: curPage}}

	for {
		issues, _, err := m.endpoint.from.Issues.ListProjectIssues(*srcproj.ID, opts)
		if err != nil {
			log.Fatal(err)
		}
		if len(issues) > 0 {
			nbIssues += len(issues)
			curPage++
			opts.Page = curPage
			for _, issue := range issues {
				switch issue.State {
				case "opened":
					nbOpened++
				case "closed":
					nbClosed++
				}
				if issue.Milestone.Title != "" {
					milestones[issue.Milestone.Title]++
				}
			}
		} else {
			break
		}
	}
	fmt.Println("OK")
	fmt.Printf("source: %d issues (%d opened, %d closed)\n", nbIssues, nbOpened, nbClosed)
	keys := make([]string, len(milestones))
	i := 0
	for k := range milestones {
		keys[i] = k
		i++
	}
	fmt.Printf("source: %d milestones: %s\n", len(keys), keys)
}
