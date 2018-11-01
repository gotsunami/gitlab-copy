package config

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/go-yaml/yaml"
	"github.com/xanzy/go-gitlab"
)

const (
	apiPath = "/api/v4"
)

type issueRange struct {
	from, to int
}

// Config contains the configuration related to source and target projects.
type Config struct {
	SrcPrj, DstPrj *project
}

// Parse reads YAML data and returns a config suitable for later
// processing.
func Parse(r io.Reader) (*Config, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}

	if err := c.SrcPrj.checkData("source"); err != nil {
		return nil, err
	}
	if err := c.DstPrj.checkData("destination"); err != nil {
		return nil, err
	}
	if err := c.checkUserTokens(); err != nil {
		return nil, err
	}
	if err := c.SrcPrj.parseIssues(); err != nil {
		return nil, err
	}
	if c.SrcPrj.LinkToTargetIssueText == "" {
		c.SrcPrj.LinkToTargetIssueText = "Closed in favor of {{.Link}}"
	}

	return c, nil
}

func (c *Config) checkUserTokens() error {
	if len(c.DstPrj.Users) == 0 {
		return nil
	}
	fmt.Printf("User tokens provided (for writing notes): %d\n", len(c.DstPrj.Users))
	fmt.Println("Checking user tokens ... ")
	for user, token := range c.DstPrj.Users {
		cl := gitlab.NewClient(nil, token)
		if err := cl.SetBaseURL(c.DstPrj.ServerURL); err != nil {
			return err
		}
		u, _, err := cl.Users.CurrentUser()
		if err != nil {
			return fmt.Errorf("Failed using the API with user '%s': %s", user, err.Error())
		}
		if u.Username != user {
			return fmt.Errorf("Token %s matches user '%s', not '%s' as defined in the config file", token, u.Username, user)
		}
	}
	fmt.Println("Tokens valid and mapping to expected users\n--")
	return nil
}
