export GOPATH=$(shell pwd)/deps:$(shell pwd)
export GOBIN=$(shell pwd)/bin

all: import search api

import:
	go get -v -d importer
	go install importer

search:
	go get -v -d search
	go install search

api:
	go get -v -d api
	go install api

goose:
	go get bitbucket.org/liamstask/goose/cmd/goose
