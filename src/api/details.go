package main

import (
	"log"
	"net/http"
	"store"
)

func detailsHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Server", "location service")
	db, err := store.NewSqlStore()
	if err != nil {
		log.Fatal(err)
	}

	ids, _ := req.URL.Query()["id"]
	p, _ := parseParams(req)

	locations, err := db.FindLocationsByIds(ids, p)
	if err != nil {
		log.Fatal(err)
	}

	err = writeLocations(locations, res, p)
	if err != nil {
		log.Fatal(err)
	}
}
