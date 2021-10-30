
PROJECT_ROOT := $(shell pwd)

.PHONY:
test:
	@go test -v -coverpkg=./... -coverprofile c.out ./...

.PHONY:
cover: test
	@go tool cover -html=c.out


