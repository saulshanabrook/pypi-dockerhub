package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/saulshanabrook/pypi-dockerhub/api"
)

var aPI = cli.Command{
	Name:  "api",
	Usage: "starts an HTTP api that returns JSON of all the releases",
	Flags: append(dbFlags, cli.StringFlag{
		Name:   "port",
		Usage:  "Port that the api is exposed",
		EnvVar: "PORT",
	}),
	Action: func(c *cli.Context) {
		http.HandleFunc("/", api.CreateHandler(getDB(c)))
		port := c.String("port")
		logrus.WithField("port", port).Info("Starting API server")
		logrus.WithError(http.ListenAndServe(":"+port, nil)).Error("Can't listen and serve")
	},
}
