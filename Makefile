ifndef VERBOSE
	MAKEFLAGS += --silent
endif

PKGS=$(shell go list ./... | grep -v /vendor)
GIT_SHA=$(shell git rev-parse --verify HEAD)
VERSION=$(shell cat VERSION)
PWD=$(shell pwd)

default: clean build

all: clean build

build: ## Make webby-bot binary
	go build -o bin/webby-bot main.go

clean: ## Cleanup build files
	rm -rf bin/*

install-tools:
	GOIMPORTS_CMD=$(shell command -v goimports 2> /dev/null)
ifndef GOIMPORTS_CMD
	go get golang.org/x/tools/cmd/goimports
endif

	GOLINT_CMD=$(shell command -v golint 2> /dev/null)
ifndef GOLINT_CMD
	go get github.com/golang/lint/golint
endif

	GODEP_CMD=$(shell command -v dep 2> /dev/null)
ifndef GODEP_CMD
	go get github.com/golang/dep/cmd/dep
endif

dep: install-tools
	dep ensure
