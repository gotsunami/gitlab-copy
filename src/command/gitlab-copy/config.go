package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/go-yaml"
)

const (
	apiPath = "/api/v3"
)

type issueRange struct {
	from, to int
}

type project struct {
	ServerURL string `yaml:"url"`
	Name      string `yaml:"project"`
	Token     string
	Issues    []string
	// Same as Issues but converted to int by parseConfig
	issues []issueRange
}

// matches checks whether issue is part of p.issues. Always
// true if p.issues is an empty list, otherwise check all entries
// and ranges, if any.
func (p *project) matches(issue int) bool {
	if len(p.issues) == 0 {
		return true
	}
	for _, i := range p.issues {
		if issue >= i.from && issue <= i.to {
			return true
		}
	}
	return false
}

// parseIssues ensure issues items are valid input data, i.e castable
// to int, ranges allowed.
func (p *project) parseIssues() error {
	p.issues = make([]issueRange, 0)
	var x [2]int
	for _, i := range p.Issues {
		vals := strings.Split(i, "-")
		if len(vals) > 2 {
			return fmt.Errorf("only one range separator allowed, '%s' not supported", vals)
		}
		if len(vals) > 1 {
			for k, p := range vals {
				num, err := strconv.ParseUint(p, 10, 64)
				if err != nil {
					return fmt.Errorf("wrong issue range in '%s': expects an integer, not '%s'", i, p)
				}
				x[k] = int(num)
			}
			if x[0] > x[1] {
				return fmt.Errorf("reverse range not allowed in '%s'", i)
			}
		} else {
			// No range
			num, err := strconv.ParseUint(vals[0], 10, 64)
			if err != nil {
				return fmt.Errorf("wrong issue value for '%s': expects an integer, not '%s'", i, vals[0])
			}
			x[0] = int(num)
			x[1] = int(num)
		}
		p.issues = append(p.issues, issueRange{from: x[0], to: x[1]})
	}
	return nil
}

type config struct {
	From, To *project
}

func checkProjectData(p *project, prefix string) error {
	if p == nil {
		return fmt.Errorf("missing %s project's data", prefix)
	}
	if p.ServerURL == "" {
		return fmt.Errorf("missing %s project's server URL", prefix)
	}
	if !strings.HasSuffix(p.ServerURL, apiPath) {
		p.ServerURL += apiPath
	}
	if p.Name == "" {
		return fmt.Errorf("missing %s project's name", prefix)
	}
	if p.Token == "" {
		return fmt.Errorf("missing %s project's token", prefix)
	}
	return nil
}

func parseConfig(name string) (*config, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	c := new(config)
	if err := yaml.Unmarshal(data, c); err != nil {
		log.Fatal(err)
	}

	if err := checkProjectData(c.From, "source"); err != nil {
		return nil, err
	}
	if err := checkProjectData(c.To, "destination"); err != nil {
		return nil, err
	}
	if err := c.From.parseIssues(); err != nil {
		return nil, err
	}

	return c, nil
}
