package main

import (
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/lookup", lookupHandler)
	http.ListenAndServe(":8080", nil)
}

type location struct {
	Id    string      `json:"id"`
	Name  string      `json:"name"`
	Shape interface{} `json:"shape,omitempty"`
}

type params struct {
	ancestors bool
	shapes    bool
}

func parseParams(req *http.Request) (params, error) {
	var p params
	p.ancestors = getBool(req, "ancestors")
	p.shapes = getBool(req, "shapes")
	return p, nil
}

func getBool(req *http.Request, key string) bool {
	val, err := strconv.ParseBool(req.URL.Query().Get(key))
	return (err == nil) && val
}
