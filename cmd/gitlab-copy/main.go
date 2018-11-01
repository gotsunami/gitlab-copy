package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/gotsunami/gitlab-copy/config"
	"github.com/gotsunami/gitlab-copy/migration"
	"github.com/gotsunami/gitlab-copy/stats"
)

func map2Human(m map[string]int) string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return strings.Join(keys, ", ")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Usage: %s [options] configfile\n", os.Args[0]))
		fmt.Fprintf(os.Stderr, `Where configfile holds YAML looks like:
from:
    url: https://gitlab.mydomain.com
    token: atoken
    project: namespace/project
    issues:
    - 5
    - 8-10
    ## Set labelsOnly to copy labels only, not issues
    # labelsOnly: true
    ## Move issues instead of copying them
    # moveIssues: true
to:
    url: https://gitlab.myotherdomain.com
    token: anothertoken
    project: namespace/project

Options:
`)
		flag.PrintDefaults()
		os.Exit(2)
	}

	apply := flag.Bool("y", false, "apply migration for real")
	version := flag.Bool("version", false, "")
	flag.Parse()

	if *version {
		fmt.Printf("Version:      %s\n", Version)
		fmt.Printf("Git revision: %s\n", GitRevision)
		fmt.Printf("Git branch:   %s\n", GitBranch)
		fmt.Printf("Go version:   %s\n", runtime.Version())
		fmt.Printf("Built:        %s\n", Built)
		fmt.Printf("OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		fmt.Fprint(os.Stderr, "Config file is missing.\n\n")
		flag.Usage()
	}
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	c, err := config.Parse(f)
	if err != nil {
		log.Fatal(err)
	}

	if !*apply {
		fmt.Println("DUMMY MODE: won't apply anything (stats only)\n--")
	}

	m, err := migration.New(c)
	if err != nil {
		log.Fatal(err)
	}
	srcproj, err := m.SourceProject(c.SrcPrj.Name)
	if err != nil {
		log.Fatal(err)
	}
	if srcproj == nil {
		log.Fatalf("source project not found on %s", c.SrcPrj.ServerURL)
	}
	fmt.Printf("source: %s at %s\n", c.SrcPrj.Name, c.SrcPrj.ServerURL)

	dstproj, err := m.DestProject(c.DstPrj.Name)
	if err != nil {
		log.Fatal(err)
	}
	if dstproj == nil {
		log.Fatalf("target project not found on %s", c.DstPrj.ServerURL)
	}
	fmt.Printf("target: %s at %s\n", c.DstPrj.Name, c.DstPrj.ServerURL)
	fmt.Println("--")

	// Find out how many issues we have
	fmt.Printf("source: finding issues ... ")

	pstats := stats.NewProject(srcproj)

	if err := pstats.ComputeStats(m.Endpoint.SrcClient); err != nil {
		log.Fatal(err)
	}
	fmt.Println("OK")
	fmt.Printf("source: %v\n", pstats)
	if len(pstats.Milestones) > 0 {
		fmt.Printf("source: %d milestone(s): %s\n", len(pstats.Milestones), map2Human(pstats.Milestones))
	}
	if len(pstats.Labels) > 0 {
		fmt.Printf("source: %d label(s): %s\n", len(pstats.Labels), map2Human(pstats.Labels))
	}

	if !c.SrcPrj.LabelsOnly {
		fmt.Printf("source: counting notes (comments), can take a while ... ")
		if err := pstats.ComputeIssueNotes(m.Endpoint.SrcClient); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\rsource: %d notes%50s\n", pstats.NbNotes, " ")
	}
	fmt.Println("--")
	if !*apply {
		if c.SrcPrj.LabelsOnly {
			fmt.Println("Will copy labels only.")
		} else {
			if c.SrcPrj.MilestonesOnly {
				fmt.Println("Will copy milestones only.")
			} else {
				action := "Copy"
				if c.SrcPrj.MoveIssues {
					action = "Move"
				}
				fmt.Printf(`Those actions will be performed:
- Copy milestones if not existing on target
- Copy all source labels on target
- %s all issues (or those specified) if not existing on target (by title)
- Copy closed status on issues, if any
- Set issue's assignee (if user exists) and milestone, if any
- Copy notes (attached to issues)
`, action)
				if c.SrcPrj.AutoCloseIssues {
					fmt.Println("- Auto-close source issues")
				}
				if c.SrcPrj.LinkToTargetIssue {
					fmt.Println("- Add a note with a link to new issue")
					fmt.Println("- Use the link text template: " + c.SrcPrj.LinkToTargetIssueText)
				}
			}
		}

		fmt.Printf("\nNow use the -y flag if that looks good to start the issues migration/label copy.\n")
		os.Exit(0)
	}

	if err := m.Migrate(); err != nil {
		log.Fatal(err)
	}
}
