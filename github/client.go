package github

import (
	"github.com/saulshanabrook/pypi-dockerhub/Godeps/_workspace/src/github.com/google/go-github/github"

	"github.com/saulshanabrook/pypi-dockerhub/Godeps/_workspace/src/golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	owner  string
	repo   string
}

func NewClient(token, owner, repo string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return &Client{client: github.NewClient(tc), owner: owner, repo: repo}
}
