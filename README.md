[![Build Status](https://travis-ci.com/paulbes/go-pedal.svg?branch=master)](https://travis-ci.com/paulbes/go-pedal)

# go-pedal
A simple client for interacting with https://developer.oslobysykkel.no/api, this demo application will print some basic information about the bike stations.

# Usage

Ensure that you have [make](https://www.gnu.org/software/make/) and [docker](https://docs.docker.com/install/) installed. You will also need a client identifier, which you can get by creating a [by sykkel account](https://developer.oslobysykkel.no/sign-up).

## Using go

Ensure that you have setup your go environment correctly: https://golang.org/doc/install

### Testing

```bash
# Display the help information
make help

# Run all the checks (linting, formatting, tests, ...)
make check
```

### Running

```bash
# As a CLI
go run cmd/pedal/main.go -client-identifier {your client identifier}

# As an API
go run cmd/api/main.go -client-identifier {your client identifier}
```

## Using docker

```bash
# Building and running the CLI
CLIENT_IDENTIFIER={your client id} make cli-build cli-run

# Building and running the API
CLIENT_IDENTIFIER={your client id} make api-build api-run
```

## As a go library

`go get github.com/paulbes/go-pedal/pedal`

## Testing the local API

```bash
# Get all stations
curl "http://localhost:8080/v1/stations"

# A specific station
curl "http://localhost:8080/v1/stations/183"
```
