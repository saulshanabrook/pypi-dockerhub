package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var deleteDockerHub = cli.Command{
	Name:  "delete-dockerHub",
	Usage: "Remove all builds in dockerHub",
	Flags: dockerHubFlags,
	Action: func(c *cli.Context) {
		dhc := getDockerHubClient(c)
		if err := dhc.DeleteAll(); err != nil {
			logrus.WithError(err).Fatal("Couldn't delete all repos in dockerHub")
		}
	},
}
