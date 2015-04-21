package main

import (
	"log"
	"net/http"
	"strconv"
)

func lookupHandler(res http.ResponseWriter, req *http.Request) {
	addHeaders(res, req)

	x, _ := strconv.ParseFloat(req.URL.Query().Get("x"), 64)
	y, _ := strconv.ParseFloat(req.URL.Query().Get("y"), 64)
	p, _ := parseParams(req)

	locations, err := (*db).FindLocationsByPoint(x, y, p)
	if err != nil {
		log.Fatal(err)
	}

	err = writeLocations(locations, res, p)
	if err != nil {
		log.Fatal(err)
	}
}
