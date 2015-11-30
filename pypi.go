package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/saulshanabrook/pypi-dockerhub/pypi"
)

func getPypiClient(c *cli.Context) *pypi.Client {
	pypiClient, err := pypi.NewClient()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Couldnt create pypi client")
	}
	return pypiClient
}
