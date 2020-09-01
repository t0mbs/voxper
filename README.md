# Voxper 0.1.1

## Introduction
> Voxper is a a REST API written in Go that acts as a vulnerability management tool, to be used in the context of security orchestration. It is intended to be middleware; ingesting, parsing and standardizing scan results from a variety of engines. This tool currently in its MVP state and should not be used in a production setting. It currently lacks basic functionality such as; encryption, API authentication keys, and user management.

> This currently only parses Snyk and Brakeman scan results, more will be added.

## Installation
> This installation assumes that the user is running the latest versions of Docker and docker-compose. Installation guides for these tools can be found [here](https://docs.docker.com/get-docker/) and [here](https://docs.docker.com/compose/install/) respectively.

```
# Clone git repository locally
git clone https://github.com/t0mbs/voxper.git

# Start docker-compose from voxper directory
cd voxper
docker-compose up

# Confirm that it is running on 8000
```

## Directory Struture
* `api` Swagger api documentation,
* `cmd` Main application for this project,
* `scripts` Non-go scripts (e.g. entrypoint, waitforit).
* `test` Folder for automated tests.

## Relevant Standards
* REST API Specs: https://restfulapi.net/
* Response Specs: https://github.com/omniti-labs/jsend
* Project layout: https://github.com/golang-standards/project-layout
* CVSS 3.1: https://www.first.org/cvss/specification-document
