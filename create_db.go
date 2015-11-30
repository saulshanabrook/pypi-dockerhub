package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/saulshanabrook/pypi-dockerhub/db"
)

var createDB = cli.Command{
	Name:  "create-db",
	Usage: "Create the database table",
	Flags: dbFlags,
	Action: func(c *cli.Context) {
		_db := getDB(c).DB.CreateTable(&db.Release{})
		if err := _db.Error; err != nil {
			logrus.WithError(err).Fatal("Couldn't drop the database table")
		}
	},
}
