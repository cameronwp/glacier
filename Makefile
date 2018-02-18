VERSION           ?= $(shell git rev-parse --short HEAD)
export

.PHONY: all build get-dependencies test mocks

get-dependencies:
	go get golang.org/x/tools/cmd/goimports
	go get github.com/golang/lint/golint
	go get github.com/kisielk/errcheck

check:
	go fmt
	go vet
	golint
	errcheck

mocks:
	GOPATH=$$(go env GOPATH) mockery -dir drain/ -all -output drain/drainmocks -outpkg drainmocks
	GOPATH=$$(go env GOPATH) mockery -dir jobqueue/ -all -output jobqueue/jobqueuemocks -outpkg jobqueuemocks

build: check
	go build

test:
	go test $$(go list ./... | grep -Ev 'migrations|scripts|vendor')
