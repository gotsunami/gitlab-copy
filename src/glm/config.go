package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-yaml"
)

type project struct {
	ServerURL string `yaml:"url"`
	Name      string `yaml:"project"`
	Token     string
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

	return c, nil
}
