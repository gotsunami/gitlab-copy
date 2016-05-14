package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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
		fmt.Printf("version: %s\n", appVersion)
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		fmt.Fprint(os.Stderr, "Config file is missing.\n\n")
		flag.Usage()
	}
	c, err := parseConfig(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	if !*apply {
		fmt.Println("DUMMY MODE: won't apply anything (stats only)\n--")
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
		fmt.Printf("source: %d milestone(s): %s\n", len(pstats.milestones), map2Human(pstats.milestones))
	}
	if len(pstats.labels) > 0 {
		fmt.Printf("source: %d label(s): %s\n", len(pstats.labels), map2Human(pstats.labels))
	}

	if !c.From.LabelsOnly {
		fmt.Printf("source: counting notes (comments), can take a while ... ")
		if err := pstats.computeIssueNotes(m.endpoint.from); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\rsource: %d notes%50s\n", pstats.nbNotes, " ")
	}
	fmt.Println("--")
	if !*apply {
		if c.From.LabelsOnly {
			fmt.Println("Will copy labels only.")
		} else {
			if c.From.MilestonesOnly {
				fmt.Println("Will copy milestones only.")
			} else {
				action := "Copy"
				if c.From.MoveIssues {
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
				if c.From.AutoCloseIssues {
					fmt.Println("- Auto-close source issues")
				}
				if c.From.LinkToTargetIssue {
					fmt.Println("- Add a note with a link to new issue")
					fmt.Println("- Use the link text template: " + c.From.LinkToTargetIssueText)
				}
			}
		}

		fmt.Printf("\nNow use the -y flag if that looks good to you to start the issues migration.\n")
		os.Exit(0)
	}

	if err := m.migrate(); err != nil {
		log.Fatal(err)
	}
}
