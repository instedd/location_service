package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/foobaz/geom"
	"github.com/foobaz/geom/encoding/geojson"
	"model"
	"net/http"
	"store"
	"strconv"
	"strings"
)

func main() {
	var port int
	var debug bool

	flag.IntVar(&port, "port", 8080, "Port where to listen for requests")
	flag.BoolVar(&debug, "debug", false, "Print debug messages")
	flag.Parse()

	if debug {
		store.SetDebug(true)
	}

	addr := fmt.Sprintf(":%d", port)

	http.HandleFunc("/lookup", lookupHandler)
	http.HandleFunc("/details", detailsHandler)
	http.HandleFunc("/children", childrenHandler)
	http.HandleFunc("/suggest", suggestHandler)
	http.ListenAndServe(addr, nil)
}

type location struct {
	Id           string      `json:"id"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	AncestorsIds []string    `json:"ancestorsIds,omitempty"`
	Ancestors    []location  `json:"ancestors,omitempty"`
	Level        int         `json:"level"`
	Lat          float64     `json:"lat"`
	Lng          float64     `json:"lng"`
	Shape        interface{} `json:"shape,omitempty"`
}

func parseParams(req *http.Request) (model.ReqOptions, error) {
	var p model.ReqOptions
	p.Ancestors = getBool(req, "ancestors")
	p.Shapes = getBool(req, "shapes")
	p.Limit = getInt(req, "limit")
	p.Offset = getInt(req, "offset")
	p.Set = req.URL.Query().Get("set")

	if len(req.URL.Query().Get("scope")) > 0 {
		p.Scope = strings.Split(req.URL.Query().Get("scope"), ",")
	} else {
		p.Scope = make([]string, 0)
	}

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

func writeLocations(locations []*model.Location, res http.ResponseWriter, p model.ReqOptions) error {
	responseLocations := make([]location, len(locations))
	for i, loc := range locations {

		point := loc.Center.T.(geom.Point)

		l := location{
			Id:           loc.Id,
			Name:         loc.Name,
			Type:         loc.TypeName,
			Level:        loc.Level,
			AncestorsIds: loc.AncestorsIds,
			Lat:          point[1],
			Lng:          point[0],
		}

		if p.Shapes {
			locationShape, err := geojson.ToGeoJSON(loc.Shape.T)
			if err != nil {
				return err
			}
			l.Shape = locationShape
		}

		if p.Ancestors {
			l.Ancestors = make([]location, len(loc.Ancestors))
			for i, anc := range loc.Ancestors {
				l.Ancestors[i] = location{
					Id:    anc.Id,
					Name:  anc.Name,
					Type:  anc.TypeName,
					Level: anc.Level,
				}
			}
		}

		responseLocations[i] = l
	}

	enc := json.NewEncoder(res)
	err := enc.Encode(responseLocations)
	if err != nil {
		return err
	}

	return nil
}
