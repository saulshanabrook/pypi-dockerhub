package dockerhub

import (
	"fmt"

	log "github.com/saulshanabrook/pypi-dockerhub/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

func (c *Client) DeleteAll() error {
	repos, err := c.allRepos()
	if err != nil {
		return err
	}
	for _, repo := range repos {
		log.WithFields(log.Fields{
			"name": repo,
		}).Info("Removing repo")
		if _, err = c.callAPI(fmt.Sprintf("v2/repositories/%v/%v/", c.dockerhubOwner, repo), "DELETE", "", 202, nil); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) allRepos() (repos []string, err error) {
	var resJSON struct {
		Results []struct {
			Name string
		}
	}
	repos = []string{}
	if _, err = c.callAPI("v2/repositories/pypi/?page=1&page_size=100", "GET", "", 200, &resJSON); err != nil {
		return
	}
	for _, res := range resJSON.Results {
		repos = append(repos, res.Name)
	}
	return
}
