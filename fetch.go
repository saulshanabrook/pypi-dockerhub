package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/saulshanabrook/pypi-dockerhub/db"
)

var fetch = cli.Command{
	Name:  "fetch",
	Usage: "Fetch packages from PyPi and add them to the database",
	Flags: append(
		[]cli.Flag{
			cli.IntFlag{
				Name:   "releases-since",
				EnvVar: "RELEASES_SINCE",
				Usage:  "Specified in seconds since epoch. Will add all package release after this time. If not specified, will get all packages",
			},
		},
		dbFlags...,
	),
	Action: func(c *cli.Context) {
		ppc := getPypiClient(c)
		var releases []db.Release
		var err error
		if time := c.Int("releases-since"); time != 0 {
			logrus.WithField("time (second since epoch)", time).Info("adding releases since")
			releases, err = ppc.ReleasesSince(int64(time))
		} else {
			logrus.Info("adding all releases")
			releases, err = ppc.AllReleases()
		}
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Fatal("couldnt get releases from pypi")
		}
		logrus.WithField("number releases", len(releases)).Info("retrieved releases")
		getDB(c).AddReleases(releases)
	},
}
