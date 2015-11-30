package dockerhub

import (
	"github.com/Sirupsen/logrus"
	"github.com/saulshanabrook/pypi-dockerhub/db"
)

func (c *Client) DeleteRepo(rel Release) error {
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
		logrus.WithFields(logrus.Fields{
			"release": repo,
		}).Info("Removing repo")
		if err = c.DeleteRepo(repo); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) allRepos() (repos []Release, err error) {
	return c.someRepos("https://hub.docker.com/v2/repositories/pypi/?page=1&page_size=100")
}

func (c *Client) someRepos(url string) (repos []Release, err error) {
	var resJSON struct {
		Next    string `json:"next"`
		Results []struct {
			Name string
		}
	}
	repos = []Release{}
	if _, err = c.callURL(url, "GET", "", 200, &resJSON); err != nil {
		return
	}
	for _, res := range resJSON.Results {
		repos = append(repos, &db.Release{Name: res.Name})
	}
	if resJSON.Next != "" {
		var moreRepos []Release
		moreRepos, err = c.someRepos(resJSON.Next)
		repos = append(repos, moreRepos...)
	}
	return repos, err
}
