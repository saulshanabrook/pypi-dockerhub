package main

import (
	"os"

	_ "github.com/lib/pq" // so we can open the databse

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "pypi-dockerHub"
	app.Usage = "Creates automated DockerHub builds for PyPi packages"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			EnvVar: "DEBUG",
		},
	}
	app.Before = setupLogging
	app.Commands = []cli.Command{
		deleteDockerHub,
		deleteGithub,
		deleteDB,
		createGithub,
		createDB,
		fetch,
		push,
		aPI,
	}
	app.Run(os.Args)
}
