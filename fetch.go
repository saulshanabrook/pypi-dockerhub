package main

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/saulshanabrook/pypi-dockerhub/db"
)

var fetch = cli.Command{
	Name:  "fetch",
	Usage: "Fetch packages from PyPi and add them to the database",
	Flags: append(
		[]cli.Flag{
			cli.BoolFlag{
				Name:   "only-new",
				EnvVar: "ONLY_NEW",
				Usage:  "Only get packages added after last fetch",
			},
		},
		dbFlags...,
	),
	Action: func(c *cli.Context) {
		ppc := getPypiClient(c)
		dbc := getDB(c)
		var releases []db.Release
		var err error
		startTime := time.Now()

		if c.Bool("only-new") {
			time, err := dbc.GetLastUpdateTime()
			if err != nil {
				logrus.WithError(err).Fatal("couldn't get last fetch time from database")
			}
			logrus.WithField("time", time).Info("adding releases since")
			releases, err = ppc.ReleasesSince(time)
		} else {
			logrus.Info("adding all releases")
			releases, err = ppc.AllReleases()
		}
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Fatal("couldnt get releases from pypi")
		}

		logrus.Info("setting last update time in DB")
		err = dbc.SetLastUpdateTime(startTime)
		if err != nil {
			logrus.WithError(err).Fatal("couldn't update last fetch time in database")
		}

		logrus.WithField("number releases", len(releases)).Info("retrieved releases")
		err = dbc.AddReleases(releases)
		if err != nil {
			logrus.WithError(err).Fatal("couldn't add all release to database")
		}
	},
}
