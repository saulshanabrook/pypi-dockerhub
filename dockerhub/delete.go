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
	return c.someRepos("https://hub.docker.com/v2/repositories/pypi/?page=1&page_size=100")
}

func (c *Client) someRepos(url string) (repos []*release.Release, err error) {
	var resJSON struct {
		Next    string `json:"next"`
		Results []struct {
			Name string
		}
	}
	repos = []*release.Release{}
	if _, err = c.callURL(url, "GET", "", 200, &resJSON); err != nil {
		return
	}
	for _, res := range resJSON.Results {
		repos = append(repos, &release.Release{res.Name, "", 0})
	}
	if resJSON.Next != "" {
		var moreRepos []*release.Release
		moreRepos, err = c.someRepos(resJSON.Next)
		repos = append(repos, moreRepos...)
	}
	return repos, err
}
