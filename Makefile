export GOPATH=$(shell pwd)/deps:$(shell pwd)
export GOBIN=$(shell pwd)/bin

DOCKERTAG := instedd/location_service
VERSION   := $(shell git describe 2>/dev/null || echo "`date -u \"+%Y%m%d.%H%M%S\"`-`git describe --always`")
PROJECTS  := importer api search

all: $(PROJECTS) goose

.PHONY: get-deps $(addsuffix -deps,$(PROJECTS))

get-deps: $(addsuffix -deps,$(PROJECTS))
	go get -v -d bitbucket.org/liamstask/goose/cmd/goose

%-deps:
	go get -v -d $*

importer:
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

docker: get-deps
	echo $(VERSION) > VERSION
	docker build --tag $(DOCKERTAG):$(VERSION) .
	docker build --tag $(DOCKERTAG) .

docker-push: docker
	docker push $(DOCKERTAG):$(VERSION)
	docker push $(DOCKERTAG)

download-ne: download-ne-0_countries download-ne-1_states_provinces

download-ne-%: importer
	mkdir -p tmp/data

	@echo "Downloading NE $*"
	curl -s http://naciscdn.org/naturalearth/10m/cultural/ne_10m_admin_$*.zip > tmp/data/ne_10m_admin_$*.zip
	unzip -qo tmp/data/ne_10m_admin_$*.zip -d tmp/data
	bin/importer -source ne tmp/data/ne_10m_admin_$*.shp
