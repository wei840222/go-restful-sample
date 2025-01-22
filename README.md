# go-restful-sample

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
