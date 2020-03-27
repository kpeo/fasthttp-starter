.PHONY: build

GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILD_TIME := $(shell date +%s)

build:
	go build -ldflags "-X main.Commit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)" .
