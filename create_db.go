package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/saulshanabrook/pypi-dockerhub/db"
)

var createDB = cli.Command{
	Name:  "create-db",
	Usage: "Create the database tables",
	Flags: dbFlags,
	Action: func(c *cli.Context) {
		logrus.Debug("Creating release table")
		_db := getDB(c).DB.CreateTable(&db.Release{})
		if err := _db.Error; err != nil {
			logrus.WithError(err).Fatal("Couldn't create the release table")
		}
		logrus.Debug("Creating lastupdatetime table")
		_db = getDB(c).DB.CreateTable(&db.LastUpdateTime{})
		if err := _db.Error; err != nil {
			logrus.WithError(err).Fatal("Couldn't create the release table")
		}
	},
}
