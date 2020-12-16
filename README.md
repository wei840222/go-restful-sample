# go-restful-sample

[![Build Status](https://cloud.drone.io/api/badges/wei840222/go-restful-sample/status.svg)](https://cloud.drone.io/wei840222/go-restful-sample)

## How to develop
### Install tools
```bash
# mockgen
go get -u github.com/golang/mock/mockgen

# swaggo
go get -u github.com/swaggo/swag/cmd/swag

# auto watch file change and hot reload server
go get -u github.com/silenceper/gowatch

# generate orm struct go code from connect to mysql
go get -u github.com/Shelnutt2/db2struct/cmd/db2struct

# setup git hooks and it will auto run go generate ./... && go test ./... on git commit
git config core.hooksPath githooks
```
### Commands
```bash
# generate go code from *.proto and mock code for testing
go generate ./...

# run testing
go test ./...

# run server, config in gowatch.yml
gowatch

# generate orm struct from mysql
db2struct --host=localhost --mysql_port=3306 --user=root --password=root1234 --gorm -d sample -t users

# build docker images
docker build -f build/dockerfile/Dockerfile .
```

## How to run manually
1. Needs all environment variables below

## Environment Variables
### Server
```bash
export DEBUG=true # enable all verborse log
export SWAGGER_BASE_URL=http://localhost:8080 # api base url on swagger ui
```
