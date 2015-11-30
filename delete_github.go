package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var deleteGithub = cli.Command{
	Name:  "delete-github",
	Usage: "Delete the github repository",
	Flags: githubFlags,
	Action: func(c *cli.Context) {
		ghc := getGithubClient(c)
		if err := ghc.DeleteRepo(); err != nil {
			logrus.WithError(err).Fatal("Couldn't delete repo in github")
		}
	},
}
