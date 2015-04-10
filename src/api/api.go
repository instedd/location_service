package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/foobaz/geom/encoding/geojson"
	"model"
	"net/http"
	"strconv"
)

func main() {
	var port int

	flag.IntVar(&port, "port", 8080, "Port where to listen for requests")
	flag.Parse()

	addr := fmt.Sprintf(":%d", port)

	http.HandleFunc("/lookup", lookupHandler)
	http.HandleFunc("/details", detailsHandler)
	http.HandleFunc("/children", childrenHandler)
	http.HandleFunc("/suggest", suggestHandler)
	http.ListenAndServe(addr, nil)
}

type location struct {
	Id        string      `json:"id"`
	Name      string      `json:"name"`
	Ancestors []string    `json:"ancestors"`
	Shape     interface{} `json:"shape,omitempty"`
}

func parseParams(req *http.Request) (model.ReqOptions, error) {
	var p model.ReqOptions
	p.Ancestors = getBool(req, "ancestors")
	p.Shapes = getBool(req, "shapes")
	p.Limit = getInt(req, "limit")
	p.Offset = getInt(req, "offset")
	p.Scope = req.URL.Query()["scope"]
	return p, nil
}

func getBool(req *http.Request, key string) bool {
	val, err := strconv.ParseBool(req.URL.Query().Get(key))
	return (err == nil) && val
}

func getInt(req *http.Request, key string) int {
	val, _ := strconv.ParseInt(req.URL.Query().Get(key), 0, 32)
	return int(val)
}

func writeLocations(locations []model.Location, res http.ResponseWriter, p model.ReqOptions) error {
	responseLocations := make([]location, len(locations))
	for i, loc := range locations {

		if p.Shapes {
			locationShape, err := geojson.ToGeoJSON(loc.Shape.T)
			if err != nil {
				return err
			}

			responseLocations[i] = location{
				Id:        loc.Id,
				Name:      loc.Name,
				Shape:     locationShape,
				Ancestors: loc.AncestorsIds,
			}
		} else {
			responseLocations[i] = location{
				Id:        loc.Id,
				Name:      loc.Name,
				Ancestors: loc.AncestorsIds,
			}
		}
	}

	enc := json.NewEncoder(res)
	err := enc.Encode(responseLocations)
	if err != nil {
		return err
	}

	return nil
}
