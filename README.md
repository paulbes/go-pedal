# go-pedal
A simple client for interacting with https://developer.oslobysykkel.no/api, this demo application will print some basic information about the bike stations.

# Usage

## As a command line tool (no go toolchain required)

1. Get your client identifier, you can get one by signing up here: https://developer.oslobysykkel.no/sign-up
2. Install docker: https://docs.docker.com/install/
3. Ensure that you have make installed: https://www.gnu.org/software/make/
4. Build the binary: `make build`
5. Run the client: `CLIENT_IDENTIFIER={your client id} make run`

```bash
$ make build
$ CLIENT_IDENTIFIER={your client id} make run
```

## As a command line tool (go toolchain required)

1. Get your client identifier, you can get one by signing up here: https://developer.oslobysykkel.no/sign-up
2. Run the client: ``

## As a go library

`go get github.com/paulbes/go-pedal/pedal`