package main

import (
	"os"

	log "github.com/saulshanabrook/pypi-dockerhub/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/saulshanabrook/pypi-dockerhub/release"
	"github.com/saulshanabrook/pypi-dockerhub/storage"

	"github.com/saulshanabrook/pypi-dockerhub/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/saulshanabrook/pypi-dockerhub/dockerhub"
	"github.com/saulshanabrook/pypi-dockerhub/github"
	"github.com/saulshanabrook/pypi-dockerhub/pypi"
)

func getDockerhubClient(c *cli.Context) *dockerhub.Client {
	dockerhubClient, err := dockerhub.NewClient(
		&dockerhub.Auth{c.GlobalString("dockerhub-username"), c.GlobalString("dockerhub-password")},
		c.GlobalString("github-owner"), c.GlobalString("github-repo"), c.GlobalString("dockerhub-owner"),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("couldnt connect to dockerhub")
	}
	return dockerhubClient
}

func getGithubClient(c *cli.Context) *github.Client {
	return github.NewClient(
		c.GlobalString("github-token"), c.GlobalString("github-owner"), c.GlobalString("github-repo"))
}

func getPypiClient(c *cli.Context) *pypi.Client {
	pypiClient, err := pypi.NewClient()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Couldnt create pypi client")
	}
	return pypiClient
}

func setupLogging(c *cli.Context) {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
		log.Debug("Set loglevel to debug")
	}
}

func getStorageClient(c *cli.Context) storage.Client {
	urlS := c.GlobalString("redis-url")
	initialTime := int64(c.GlobalInt("initial-time"))
	if urlS == "" {
		return storage.NewTestingClient(initialTime)
	}
	stor, err := storage.NewRedisClient(c.GlobalString("redis-url"), initialTime)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Couldnt create redis client")
	}
	return stor
}

func main() {
	app := cli.NewApp()
	app.Name = "pypi-dockerhub"
	app.Usage = "Creates automed Dockerhub builds for all PyPi packages. It will process all releases since the last date it recorded."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "github-owner",
			EnvVar: "GITHUB_OWNER",
		},
		cli.StringFlag{
			Name:   "github-repo",
			EnvVar: "GITHUB_REPO",
		},
		cli.StringFlag{
			Name:   "github-token",
			EnvVar: "GITHUB_TOKEN",
		},
		cli.StringFlag{
			Name:   "dockerhub-username",
			EnvVar: "DOCKERHUB_USERNAME",
		},
		cli.StringFlag{
			Name:   "dockerhub-password",
			EnvVar: "DOCKERHUB_PASSWORD",
		},
		cli.StringFlag{
			Name:   "dockerhub-owner",
			EnvVar: "DOCKERHUB_OWNER",
		},
		cli.StringFlag{
			Name:   "redis-url",
			Usage:  "if not provided, then will not persist the last update time, and you must provide `initial-time`",
			EnvVar: "REDIS_URL",
		},
		cli.IntFlag{
			Name:   "initial-time",
			Usage:  "If provided, this time (in seconds since epoch) will overwrite the recorded last update time",
			EnvVar: "INITIAL_DATE",
		},
		cli.BoolFlag{
			Name:   "debug",
			EnvVar: "DEBUG",
		},
		cli.StringFlag{
			Name:   "test-name",
			Usage:  "If provided, will not query pypi for all packages, instead just use this name",
			EnvVar: "TEST_NAME",
		},
		cli.StringFlag{
			Name:   "test-version",
			Usage:  "If provided, will not query pypi for all packages, instead just use this version",
			EnvVar: "TEST_VERSION",
		},
	}

	app.Commands = []cli.Command{{
		Name:  "remove-dockerhub",
		Usage: "Remove all builds in dockerhub",
		Action: func(c *cli.Context) {
			setupLogging(c)
			dockerhubContext := getDockerhubClient(c)
			if err := dockerhubContext.DeleteAll(); err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Fatal("Couldnt delete all repos in dockerhub")
			}
		},
	}}

	app.Name = "sync"
	app.Usage = "Create automated dockerhub builds for pypi packages"
	app.Action = func(c *cli.Context) {
		setupLogging(c)
		var rels []*release.Release
		var storageClient storage.Client
		if c.GlobalIsSet("test-name") {
			rels = []*release.Release{{
				Name:    c.GlobalString("test-name"),
				Version: c.GlobalString("test-version"),
				Time:    0,
			}}
		} else {
			pypiClient := getPypiClient(c)
			storageClient = getStorageClient(c)
			time, err := storageClient.GetTime()
			if err != nil {
				log.WithFields(log.Fields{
					"err":     err,
					"storage": storageClient,
				}).Fatal("couldnt get initial time")
			}
			rels, err = pypiClient.ReleasesSince(time)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Fatal("couldnt get releases from pypi")
			}
		}
		dockerhubClient := getDockerhubClient(c)
		githubClient := getGithubClient(c)

		for _, rel := range rels {
			rel.Log().Info("Adding Release")

			if err := githubClient.AddRelease(rel); err != nil {
				rel.Log().WithFields(log.Fields{
					"err": err,
				}).Fatal("couldnt add releases to github")
			}

			if err := dockerhubClient.AddRelease(rel); err != nil {
				rel.Log().WithFields(log.Fields{
					"err": err,
				}).Fatal("couldnt add release to dockerhub")
			}
			if storageClient != nil {
				if err := storageClient.SetTime(rel.Time); err != nil {
					rel.Log().WithFields(log.Fields{
						"err":     err,
						"storage": storageClient,
					}).Fatal("couldnt set time in storage client")
				}
			}
		}
	}
	app.Run(os.Args)

}
