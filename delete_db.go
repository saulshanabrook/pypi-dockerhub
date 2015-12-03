package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/saulshanabrook/pypi-dockerhub/db"
)

var deleteDB = cli.Command{
	Name:  "delete-db",
	Usage: "Drop the database table",
	Flags: dbFlags,
	Action: func(c *cli.Context) {
		logrus.Debug("Dropping release table")
		_db := getDB(c).DB.DropTable(&db.Release{})
		if err := _db.Error; err != nil {
			logrus.WithError(err).Fatal("Couldn't drop the release table")
		}

		logrus.Debug("Dropping lastupdatetime table")
		_db = getDB(c).DB.DropTable(&db.LastUpdateTime{})
		if err := _db.Error; err != nil {
			logrus.WithError(err).Fatal("Couldn't drop the lastupdatetime table")
		}
	},
}
