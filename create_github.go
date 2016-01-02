package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var createGithub = cli.Command{
	Name:  "create-github",
	Usage: "Create the the github repository",
	Flags: githubFlags,
	Action: func(c *cli.Context) {
		ghc := getGithubClient(c)
		if err := ghc.CreateRepo(); err != nil {
			logrus.WithError(err).Fatal("Couldn't create repo in github")
		}
	},
}
