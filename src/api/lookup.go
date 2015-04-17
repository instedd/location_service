package main

import (
	"log"
	"net/http"
	"store"
	"strconv"
)

func lookupHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Server", "location service")
	res.Header().Add("Access-Control-Allow-Origin", "*")
	db, err := store.NewSqlStore()
	if err != nil {
		log.Fatal(err)
	}

	x, _ := strconv.ParseFloat(req.URL.Query().Get("x"), 64)
	y, _ := strconv.ParseFloat(req.URL.Query().Get("y"), 64)
	p, _ := parseParams(req)

	locations, err := db.FindLocationsByPoint(x, y, p)
	if err != nil {
		log.Fatal(err)
	}

	err = writeLocations(locations, res, p)
	if err != nil {
		log.Fatal(err)
	}
}
