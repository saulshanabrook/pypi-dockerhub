# pypi-dockerhub

I recently saw the cool looking [gping](https://github.com/orf/gping)
python package. I wanted to try it out, but I have stopped trying to use
Python on my host mac system and instead run it always through Docker.

So I created a `Dockerfile` like this in a temporary directory:
```
FROM python:3

RUN pip install pinggraph
```

And then built it:

```
$ docker build -t gping .
```

And finally could use it:

```
$ docker run --rm gping gping google.com
```

What if instead I could just do:

```
$ docker run --rm -it pypi/pinggraph gping google.com
```

That is now possible.

## Structure

### Github

The  `saulshanabrook/pypi-dockerhub_` is built through automated builds on Docker Hub.

Each subdirectory of the `saulshanabrook/pypi-dockerhub_` corresponds
to a package name on pypi. Inside of each is just a Dockerfile that `pip` installs
that package. All Dockerfiles extend from python 3.

There is only one branch (`master`).

Tags are added for each release of each package in the format `<NAME>@<VERSION>` so like `pinggraph@0.0.9`.

### Docker Hub

Each pypi package gets its own docker image under the `pypi` organization.

There is an automated build for each tag, with the subdirectory as the name of the package
and the docker tag equal to the pypi version. Also there is a `latest` tag pointing to the
master branch for that subdirectory.

## Running

```
$ go install github.com/saulshanabrook/pypi-dockerhub
$ pypi-dockerhub --help
NAME:
   sync - Create automated dockerhub builds for pypi packages

USAGE:
   sync [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   remove-dockerhub	Remove all builds in dockerhub
   help, h		Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --github-owner 		 [$GITHUB_OWNER]
   --github-repo 		 [$GITHUB_REPO]
   --github-token 		 [$GITHUB_TOKEN]
   --dockerhub-username 	 [$DOCKERHUB_USERNAME]
   --dockerhub-password 	 [$DOCKERHUB_PASSWORD]
   --dockerhub-owner 		 [$DOCKERHUB_OWNER]
   --redis-url 			if not provided, then will not persist the last update time, and you must provide `initial-time` [$REDIS_URL]
   --initial-time "0"		If provided, this time (in seconds since epoch) will overwrite the recorded last update time [$INITIAL_DATE]
   --debug			 [$DEBUG]
   --help, -h			show help
   --version, -v		print the version
```

## Deploying
1. Deploy this to Heroku [![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)
2. Do an initial run, starting at some time (seconds from epoch) `heroku run pypi-dockerhub --initial-time 1445304164`
3. Add scheduler task to run `pypi-dockerhub` every couple of hours.

