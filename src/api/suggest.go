package main

import (
	"log"
	"net/http"
)

func suggestHandler(res http.ResponseWriter, req *http.Request) {
	addHeaders(res, req)

	name := req.URL.Query().Get("name")
	p, _ := parseParams(req)

	locations, err := (*db).FindLocationsByName(name, p)
	if err != nil {
		log.Fatal(err)
	}

	err = writeLocations(locations, res, p)
	if err != nil {
		log.Fatal(err)
	}
}
