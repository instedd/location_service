package main

import (
	"log"
	"net/http"
	"store"
)

func suggestHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Server", "location service")
	db, err := store.NewSqlStore()
	if err != nil {
		log.Fatal(err)
	}

	name := req.URL.Query().Get("name")
	p, _ := parseParams(req)

	locations, err := db.FindLocationsByName(name, p)
	if err != nil {
		log.Fatal(err)
	}

	err = writeLocations(locations, res, p)
	if err != nil {
		log.Fatal(err)
	}
}
