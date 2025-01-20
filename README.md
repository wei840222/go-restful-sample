# go-restful-sample

[![Build Status](https://cloud.drone.io/api/badges/wei840222/go-restful-sample/status.svg)](https://cloud.drone.io/wei840222/go-restful-sample)

## How to develop
### Install tools
```bash
# mockgen
go install go.uber.org/mock/mockgen@latest

# swaggo
go install github.com/swaggo/swag/cmd/swag@latest

# auto watch file change and hot reload server
go install github.com/silenceper/gowatch@latest

```
### Commands
```bash
# generate go code from *.proto and mock code for testing
go generate ./...

# run testing
go test ./...

# run server, config in gowatch.yml
gowatch

```

## How to run manually
1. Needs all environment variables below

## Environment Variables
### Server
```bash
export DEBUG=true # enable all verborse log
export SWAGGER_BASE_URL=http://localhost:8080 # api base url on swagger ui
```
