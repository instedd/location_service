package main

import (
	"log"
	"net/http"
	"store"
)

func childrenHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Server", "location service")
	res.Header().Add("Access-Control-Allow-Origin", "*")
	db, err := store.NewSqlStore()
	if err != nil {
		log.Fatal(err)
	}

	parentId := req.URL.Query().Get("id")
	p, _ := parseParams(req)

	locations, err := db.FindLocationsByParent(parentId, p)
	if err != nil {
		log.Fatal(err)
	}

	err = writeLocations(locations, res, p)
	if err != nil {
		log.Fatal(err)
	}
}
