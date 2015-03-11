export GOPATH=$(shell pwd)/deps:$(shell pwd)

import:
	go get -v -d import
	go install import
