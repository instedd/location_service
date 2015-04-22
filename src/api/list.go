package main

import (
	"log"
	"net/http"
)

func listHandler(res http.ResponseWriter, req *http.Request) {
	addHeaders(res, req)

	p, _ := parseParams(req)
	locations, err := (*db).FindLocations(p)
	if err != nil {
		log.Fatal(err)
	}

	err = writeLocations(locations, res, p)
	if err != nil {
		log.Fatal(err)
	}
}
