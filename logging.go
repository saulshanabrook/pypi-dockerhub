package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func setupLogging(c *cli.Context) error {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
		log.Debug("Set loglevel to debug")
	}
	return nil
}
