package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/saulshanabrook/pypi-dockerhub/dockerhub"
)

func getDockerHubClient(c *cli.Context) *dockerhub.Client {
	dockerHubClient, err := dockerhub.NewClient(
		&dockerhub.Auth{Username: c.String("dockerhub-username"), Password: c.String("dockerhub-password")},
		&dockerhub.Repo{Owner: c.String("dockerhub-owner"), Name: c.String("dockerhub-username")},
		getGithubRepo(c),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("couldnt connect to dockerHub")
	}
	return dockerHubClient
}

var dockerHubFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "dockerhub-username",
		EnvVar: "DOCKERHUB_USERNAME",
	},
	cli.StringFlag{
		Name:   "dockerhub-password",
		EnvVar: "DOCKERHUB_PASSWORD",
	},
	cli.StringFlag{
		Name:   "dockerhub-owner",
		EnvVar: "DOCKERHUB_OWNER",
	},
}
