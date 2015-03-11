package store

import (
	"github.com/foobaz/geom/encoding/geojson"
	elastigo "github.com/mattbaird/elastigo/lib"
	"model"
)

type esStore struct {
	conn *elastigo.Conn
}

type esLocation struct {
	Name  string      `json:"name"`
	Shape interface{} `json:"shape"`
}

func NewElasticSearchStore() (Store, error) {
	conn := elastigo.NewConn()
	conn.CreateIndex("location_service")

	type hash map[string]interface{}
	err := conn.PutMapping("location_service", "location", model.Location{}, elastigo.MappingOptions{
		Properties: hash{"shape": hash{
			"type": "geo_shape",
		}},
	})

	if err != nil {
		return nil, err
	}

	return esStore{conn}, nil
}

func (self esStore) AddLocation(location *model.Location) error {
	shape, err := geojson.ToGeoJSON(location.Shape)
	if err != nil {
		return err
	}

	_, err = self.conn.Index("location_service", "location", "", nil, esLocation{
		Name:  location.Name,
		Shape: shape,
	})
	return err
}

func (self esStore) Begin() Store {
	return self
}

func (self esStore) Flush() {
	self.conn.Flush()
}
