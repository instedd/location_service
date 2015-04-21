package main

import (
	"log"
	"net/http"
)

func suggestHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Server", "location service")
	res.Header().Add("Access-Control-Allow-Origin", "*")
	db, err := store.NewSqlStore()
	if err != nil {
		log.Fatal(err)
	}

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
