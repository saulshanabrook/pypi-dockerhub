package github

import "github.com/google/go-github/github"

func (c *Client) DeleteRepo() error {
	_, err := c.client.Repositories.Delete(c.Owner, c.Name)
	return err
}

func (c *Client) CreateRepo() error {
	repo := github.Repository{
		Name: &c.Name,
	}
	_, _, err := c.client.Repositories.Create(c.Owner, &repo)
	return err
}
