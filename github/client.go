package github

import (
	"github.com/google/go-github/github"

	"golang.org/x/oauth2"
)

type Repo struct {
	Owner string
	Name  string
}

type Client struct {
	*Repo
	client *github.Client
}

func NewClient(token string, repo *Repo) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return &Client{repo, github.NewClient(tc)}
}
