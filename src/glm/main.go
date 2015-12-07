package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func map2human(m map[string]int) string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return strings.Join(keys, ", ")
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
