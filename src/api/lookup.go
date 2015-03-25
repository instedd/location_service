package main

import (
	"encoding/json"
	"github.com/foobaz/geom/encoding/geojson"
	"log"
	"net/http"
	"store"
	"strconv"
)

func lookupHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Server", "location service")
	db, err := store.NewSqlStore()
	if err != nil {
		log.Fatal(err)
	}

	x, _ := strconv.ParseFloat(req.URL.Query().Get("x"), 64)
	y, _ := strconv.ParseFloat(req.URL.Query().Get("y"), 64)
	p, _ := parseParams(req)

	locations, err := db.FindLocationsByPoint(x, y, p.shapes)

	if err != nil {
		log.Fatal(err)
	}

	responseLocations := make([]location, len(locations))

	for i, loc := range locations {

		if p.shapes {
			locationShape, err := geojson.ToGeoJSON(loc.Shape)
			if err != nil {
				log.Fatal(err)
			}

			responseLocations[i] = location{
				Id:    loc.Id,
				Name:  loc.Name,
				Shape: locationShape,
			}
		} else {
			responseLocations[i] = location{
				Id:   loc.Id,
				Name: loc.Name,
			}
		}
	}

	enc := json.NewEncoder(res)
	err = enc.Encode(responseLocations)
	if err != nil {
		log.Fatal(err)
	}
}
