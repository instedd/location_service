package main

import (
	"log"
	"net/http"
)

func childrenHandler(res http.ResponseWriter, req *http.Request) {
	addHeaders(res, req)

	parentId := req.URL.Query().Get("id")
	p, _ := parseParams(req)

	locations, err := (*db).FindLocationsByParent(parentId, p)
	if err != nil {
		log.Fatal(err)
	}

	err = writeLocations(locations, res, p)
	if err != nil {
		log.Fatal(err)
	}
}
