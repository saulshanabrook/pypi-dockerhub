package dockerhub

import (
	log "github.com/saulshanabrook/pypi-dockerhub/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/saulshanabrook/pypi-dockerhub/release"
)

func (c *Client) DeleteRepo(rel *release.Release) error {
	res, err := c.callRepo(rel, "", "DELETE", "", 202, nil)
	if err != nil {
		return err
	}
	return res.Body.Close()
}

func (c *Client) DeleteAll() error {
	repos, err := c.allRepos()
	if err != nil {
		return err
	}
	for _, repo := range repos {
		log.WithFields(log.Fields{
			"name": repo,
		}).Info("Removing repo")
		if err = c.DeleteRepo(repo); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) allRepos() (repos []*release.Release, err error) {
	var resJSON struct {
		Results []struct {
			Name string
		}
	}
	repos = []*release.Release{}
	if _, err = c.callAPI("v2/repositories/pypi/?page=1&page_size=100", "GET", "", 200, &resJSON); err != nil {
		return
	}
	for _, res := range resJSON.Results {
		repos = append(repos, &release.Release{res.Name, "", 0})
	}
	return
}
