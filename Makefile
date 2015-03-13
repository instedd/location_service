export GOPATH=$(shell pwd)/deps:$(shell pwd)
export GOBIN=$(shell pwd)/bin

import:
	go get -v -d importer
	go install importer

goose:
	go get bitbucket.org/liamstask/goose/cmd/goose
