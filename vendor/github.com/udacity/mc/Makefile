.PHONY: all check clean dep dev format install release tag test

GOPATH_APP_FILES := $(shell go list ./... | grep -v /vendor/)
REL_APP_FILES := $(shell find . -type f | grep -vE '/(vendor|dist)/' | grep '\.go')
GO_VERSION ?= $(shell go version | cut -d ' ' -f3)
NAME := mc
ORG := udacity
VERSION ?= $(shell git describe --abbrev=0 --tags)


all: install

check:
	IMPORT_ISSUES=$(shell goimports -d $(REL_APP_FILES))
	if [ "$(IMPORT_ISSUES)" ]; then \
		exit 1; \
	fi
	FORMAT_ISSUES=$(shell gofmt -d $(REL_APP_FILES))
	if [ "$(FORMAT_ISSUES)" ]; then \
		exit 1; \
	fi
	go vet $(GOPATH_APP_FILES)
	golint -set_exit_status $(GOPATH_APP_FILES)
	errcheck $(GOPATH_APP_FILES)

clean:
	rm -f $(PWD)/dist/*

dep:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

dev:
	go get -u github.com/kisielk/errcheck
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/lint/golint

dist: clean
	@mkdir -p $(PWD)/dist
	GOOS=darwin GOARCH=amd64 go build -o dist/$(NAME)-$(VERSION).darwin-amd64.$(GO_VERSION)
	GOOS=linux GOARCH=amd64 go build -o dist/$(NAME)-$(VERSION).linux-amd64.$(GO_VERSION)

format:
	goimports -w .
	go fmt $(GOPATH_APP_FILES)
	go vet $(GOPATH_APP_FILES)
	golint $(GOPATH_APP_FILES)
	errcheck $(GOPATH_APP_FILES)

install: dev dep
	go install

release: dist
	bin/release

tag:
	bin/tag

test: check
	go test $(GOPATH_APP_FILES)
