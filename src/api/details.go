package main

import (
	"log"
	"net/http"
	"strings"
)

func detailsHandler(res http.ResponseWriter, req *http.Request) {
	addHeaders(res, req)

	ids := strings.Split(req.URL.Query().Get("id"), ",")
	p, _ := parseParams(req)

	locations, err := (*db).FindLocationsByIds(ids, p)
	if err != nil {
		log.Fatal(err)
	}

	err = writeLocations(locations, res, p)
	if err != nil {
		log.Fatal(err)
	}
}
