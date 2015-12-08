package main

import (
	"errors"
	"fmt"
	"net/http"

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
	params                 *config
	endpoint               *endpoint
	srcProject, dstProject *gitlab.Project
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
	m := &migration{params: c, endpoint: &endpoint{fromgl, togl}}
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
	p, err := m.project(m.endpoint.from, name)
	if err != nil {
		return nil, err
	}
	m.srcProject = p
	return p, nil
}

func (m *migration) destProject(name string) (*gitlab.Project, error) {
	p, err := m.project(m.endpoint.to, name)
	if err != nil {
		return nil, err
	}
	m.dstProject = p
	return p, nil
}

// Performs the issues migration.
func (m *migration) migrate() error {
	if m.srcProject == nil || m.dstProject == nil {
		return errors.New("nil project.")
	}
	fmt.Println("Migrating ...")

	curPage := 1
	opts := &gitlab.ListProjectIssuesOptions{ListOptions: gitlab.ListOptions{PerPage: resultsPerPage, Page: curPage}}
	issues, _, err := m.endpoint.from.Issues.ListProjectIssues(m.srcProject.ID, opts)
	if err != nil {
		return err
	}
	if len(issues) > 0 {
	}
	return nil
}
