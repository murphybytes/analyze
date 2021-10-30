
PROJECT_ROOT := $(shell pwd)
GOPATH ?= $(shell go env GOPATH)

.PHONY:
test:
	@go test -v -coverpkg=./... -coverprofile c.out ./...

.PHONY:
cover: test
	@go tool cover -html=c.out

$(GOPATH)/bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.39.0

lint: $(GOPATH)/bin/golangci-lint
	golangci-lint run 

