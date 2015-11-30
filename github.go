package main

import (
	"github.com/codegangsta/cli"
	"github.com/saulshanabrook/pypi-dockerhub/github"
)

func getGithubRepo(c *cli.Context) *github.Repo {
	return &github.Repo{c.String("github-owner"), c.String("github-repo")}
}

func getGithubClient(c *cli.Context) *github.Client {
	return github.NewClient(
		c.String("github-token"), getGithubRepo(c))
}

var githubFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "github-owner",
		EnvVar: "GITHUB_OWNER",
	},
	cli.StringFlag{
		Name:   "github-repo",
		EnvVar: "GITHUB_REPO",
	},
	cli.StringFlag{
		Name:   "github-token",
		EnvVar: "GITHUB_TOKEN",
	},
}
