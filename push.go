package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/saulshanabrook/pypi-dockerhub/db"
	"github.com/saulshanabrook/pypi-dockerhub/dockerhub"
	"github.com/saulshanabrook/pypi-dockerhub/github"
)

func pushRelease(dbc *db.Client, dhc *dockerhub.Client, ghc *github.Client, rel *db.Release) error {
	logrus.WithFields(rel.Fields()).Info("Processing Release")
	if !rel.AddedGithub {
		logrus.WithFields(rel.Fields()).Debug("Adding to Github")
		if err := ghc.AddRelease(rel); err != nil {
			return err
		}
		rel.AddedGithub = true
		logrus.WithFields(rel.Fields()).Debug("Recording added to Github")
		if _db := dbc.DB.Save(rel); _db.Error != nil {
			return _db.Error
		}
	}
	if !rel.AddedDockerHub {
		logrus.WithFields(rel.Fields()).Debug("Adding to DockerHub")
		if err := dhc.AddRelease(rel); err != nil {
			return err
		}
		rel.AddedDockerHub = true
		logrus.WithFields(rel.Fields()).Debug("Recording added to DockerHub")
		if _db := dbc.DB.Save(rel); _db.Error != nil {
			return _db.Error
		}
	}
	if !rel.TriggeredDockerHub {
		logrus.WithFields(rel.Fields()).Debug("Triggering on DockerHub")

		if err := dhc.TriggerRelease(rel); err != nil {
			return err
		}
		logrus.WithFields(rel.Fields()).Debug("Recording triggering on DockerHub")
		rel.TriggeredDockerHub = true
		if _db := dbc.DB.Save(rel); _db.Error != nil {
			return _db.Error
		}
	}
	return nil
}

var push = cli.Command{
	Name:  "push",
	Usage: "Add PyPi release from database to Github and Dockerhub, so they are built",
	Flags: append(
		[]cli.Flag{
			cli.IntFlag{
				Name:   "releases-since",
				EnvVar: "RELEASES_SINCE",
				Usage:  "Specified in seconds since epoch. Will add all package release after this time.",
			},
		},
		dbFlags...,
	),
	Action: func(c *cli.Context) {
		dbc := getDB(c)
		rels, err := dbc.GetToProcess()
		if err != nil {
			logrus.WithError(err).Fatal("Couldn't get all releases from the database")
		}
		dhc := getDockerHubClient(c)
		ghc := getGithubClient(c)
		for _, rel := range rels {
			if err := pushRelease(dbc, dhc, ghc, &rel); err != nil {
				logrus.WithError(err).WithFields(rel.Fields()).Fatal("Couldn't delete all repos in dockerHub")
			}
		}
	},
}
