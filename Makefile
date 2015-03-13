export GOPATH=$(shell pwd)/deps:$(shell pwd)
export GOBIN=$(shell pwd)/bin

all: import search

import:
	go get -v -d importer
	go install importer

search:
	go get -v -d search
	go install search

goose:
	go get bitbucket.org/liamstask/goose/cmd/goose
