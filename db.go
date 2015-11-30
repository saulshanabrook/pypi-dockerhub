package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/saulshanabrook/pypi-dockerhub/db"
)

func getDB(c *cli.Context) *db.Client {
	pgURL := c.String("database-url")
	log.WithField("url", pgURL).Debug("Opening connection to database")
	client, err := db.NewClient(pgURL)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Couldnt create database client")
	}
	return client
}

var dbFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "database-url",
		Usage:  "Postgres database URL",
		EnvVar: "DATABASE_URL",
	},
}
