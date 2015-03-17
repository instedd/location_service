package main

import (
	"encoding/json"
	"github.com/foobaz/geom/encoding/geojson"
	"log"
	"net/http"
	"store"
	"strconv"
)

func main() {
	http.HandleFunc("/lookup", handler)
	http.ListenAndServe(":8080", nil)
}

type location struct {
	Id    string      `json:"id"`
	Name  string      `json:"name"`
	Shape interface{} `json:"shape,omitempty"`
}

func handler(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Server", "location service")
	db, err := store.NewSqlStore()
	if err != nil {
		log.Fatal(err)
	}

	x, _ := strconv.ParseFloat(req.URL.Query().Get("x"), 64)
	y, _ := strconv.ParseFloat(req.URL.Query().Get("y"), 64)

	locations, err := db.FindLocationsByPoint(x, y, true)
	if err != nil {
		log.Fatal(err)
	}

	responseLocations := make([]location, len(locations))

	for i, loc := range locations {

		locationShape, err := geojson.ToGeoJSON(loc.Shape)
		if err != nil {
			log.Fatal(err)
		}

		responseLocations[i] = location{
			Id:   loc.Id,
			Name: loc.Name,
			Shape: locationShape,
		}
	}

	enc := json.NewEncoder(res)
	err = enc.Encode(responseLocations)
	if err != nil {
		log.Fatal(err)
	}
}
