# pypi-dockerHub

I recently saw the cool looking [gping](https://github.com/orf/gping)
python package. I wanted to try it out, but I have stopped trying to use
Python on my host system and instead try to run everything through Docker.

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

..but I wish I could just:

```
$ docker run --rm -it pypi/pinggraph gping google.com
```

That is now possible.

## Cavaets

* All packages are installed on top of the `python:3` image.
* No system dependencies are installed.

## Structure

### Github

The  `saulshanabrook/pypi-dockerHub_` is built through automated builds on Docker Hub.

Each subdirectory of the `saulshanabrook/pypi-dockerHub_` corresponds
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
$ pypi-dockerHub --help
NAME:
   sync - Create automated dockerHub builds for pypi packages

USAGE:
   sync [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   remove-dockerHub	Remove all builds in dockerHub
   help, h		Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --github-owner 		 [$GITHUB_OWNER]
   --github-repo 		 [$GITHUB_REPO]
   --github-token 		 [$GITHUB_TOKEN]
   --dockerHub-username 	 [$DOCKERHUB_USERNAME]
   --dockerHub-password 	 [$DOCKERHUB_PASSWORD]
   --dockerHub-owner 		 [$DOCKERHUB_OWNER]
   --redis-url 			if not provided, then will not persist the last update time, and you must provide `initial-time` [$REDIS_URL]
   --initial-time "0"		If provided, this time (in seconds since epoch) will overwrite the recorded last update time [$INITIAL_DATE]
   --debug			 [$DEBUG]
   --help, -h			show help
   --version, -v		print the version
```


## Development

Initially:

```bash
echo 'DOCKERHUB_OWNER=...
DOCKERHUB_PASSWORD=...
DOCKERHUB_USERNAME=...
GITHUB_OWNER=...
GITHUB_REPO=...
GITHUB_TOKEN=...
' > .env

docker-compose --x-networking up -d db
docker-compose --x-networking run app go run *.go --debug create-db
docker-compose --x-networking up -d app
docker-compose --x-networking run app go run *.go --debug fetch
docker-compose --x-networking run app go run *.go --debug create-github
docker-compose --x-networking run app go run *.go --debug push
```

To query app:

```bash
open http://$(docker-machine ip default):8000/
```

To process updates:

```bash
docker-compose --x-networking up -d db
docker-compose --x-networking run app go run *.go --debug fetch --only-new
docker-compose --x-networking run app go run *.go --debug push
```

## Deploying
1. Deploy this to Heroku [![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)
2. Do an initial run, starting at some time (seconds from epoch) `heroku run pypi-dockerHub --initial-time 1445304164`
3. Add scheduler task to run `pypi-dockerHub` every couple of hours.

