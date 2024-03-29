package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/foobaz/geom"
	"github.com/foobaz/geom/encoding/geojson"
	"log"
	"model"
	"net/http"
	"store"
	"strconv"
	"strings"
)

var db *store.Store

func main() {
	var port int
	var debug bool

	flag.IntVar(&port, "port", 8080, "Port where to listen for requests")
	flag.BoolVar(&debug, "debug", false, "Print debug messages")
	flag.Parse()

	if debug {
		store.SetDebug(true)
	}

	store, err := store.NewSqlStore()
	if err != nil {
		panic(err)
	}

	db = &store
	addr := fmt.Sprintf(":%d", port)

	http.HandleFunc("/lookup", lookupHandler)
	http.HandleFunc("/details", detailsHandler)
	http.HandleFunc("/children", childrenHandler)
	http.HandleFunc("/suggest", suggestHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/", emptyResponseHandler)

	fmt.Printf("Starting on http://0.0.0.0:%d/\n", port)
	http.ListenAndServe(addr, httpLog(http.DefaultServeMux))
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
	Set          string      `json:"set"`
	Shape        interface{} `json:"shape,omitempty"`
}

func httpLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func parseParams(req *http.Request) (model.ReqOptions, error) {
	var p model.ReqOptions
	p.Ancestors = getBool(req, "ancestors")
	p.Shapes = getBool(req, "shapes")
	p.Simplify = getFloat(req, "simplify")
	p.Limit = getInt(req, "limit")
	p.Offset = getInt(req, "offset")
	p.Set = req.URL.Query().Get("set")
	p.Object = req.URL.Query().Get("object")

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

func getFloat(req *http.Request, key string) float32 {
	val, _ := strconv.ParseFloat(req.URL.Query().Get(key), 32)
	return float32(val)
}

func addHeaders(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Server", "location service")
	res.Header().Add("Access-Control-Allow-Origin", "*")
	res.Header().Add("Access-Control-Allow-Headers", req.Header.Get("Access-Control-Request-Headers"))
}

func writeLocations(locations []*model.Location, res http.ResponseWriter, p model.ReqOptions) error {
	var err error
	enc := json.NewEncoder(res)

	if len(p.Object) > 0 {
		responseLocations := make(map[string](*location))
		for _, loc := range locations {
			responseLocations[loc.Id] = buildLocation(loc, p)
		}
		responseObject := make(map[string](*map[string](*location)))
		responseObject[p.Object] = &responseLocations
		err = enc.Encode(responseObject)
	} else {
		responseLocations := make([]*location, len(locations))
		for i, loc := range locations {
			responseLocations[i] = buildLocation(loc, p)
		}
		err = enc.Encode(responseLocations)
	}

	if err != nil {
		return err
	}
	return nil
}

func buildLocation(loc *model.Location, p model.ReqOptions) *location {
	point := loc.Center.T.(geom.Point)
	set := strings.Split(loc.Id, ":")[0]

	l := location{
		Id:           loc.Id,
		Name:         loc.Name,
		Type:         loc.TypeName,
		Level:        loc.Level,
		AncestorsIds: loc.AncestorsIds,
		Set:          set,
		Lat:          point[1],
		Lng:          point[0],
	}

	if p.Shapes {
		locationShape, err := geojson.ToGeoJSON(loc.Shape.T)
		if err == nil {
			l.Shape = locationShape
		}
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

	return &l
}

func emptyResponseHandler(res http.ResponseWriter, req *http.Request) {
	addHeaders(res, req)

	enc := json.NewEncoder(res)

	err := enc.Encode(map[string]string{"status": "200 OK"})
	
	if err != nil {
		log.Fatal(err)
	}
}
