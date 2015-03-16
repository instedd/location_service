package store

import (
	"encoding/json"
	"github.com/foobaz/geom/encoding/geojson"
	elastigo "github.com/mattbaird/elastigo/lib"
	"model"
)

type esStore struct {
	conn *elastigo.Conn
}

type esLocation struct {
	ParentId *string `json:"parent_id"`
	Name     string  `json:"name"`
}

type esGeometry struct {
	Shape interface{} `json:"shape"`
}

type hash map[string]interface{}

func NewElasticSearchStore() (Store, error) {
	conn := elastigo.NewConn()
	conn.CreateIndex("location_service")

	err := conn.PutMapping("location_service", "location", esLocation{}, elastigo.MappingOptions{})
	if err != nil {
		return nil, err
	}

	err = conn.PutMapping("location_service", "geometry", esGeometry{}, elastigo.MappingOptions{
		Parent: &elastigo.ParentOptions{Type: "location"},
		Properties: hash{
			"shape": hash{
				"type": "geo_shape",
				"tree": "quadtree",
			},
		},
	})

	return esStore{conn}, nil
}

func (self esStore) AddLocation(location *model.Location) error {
	var err error
	shape, err := geojson.ToGeoJSON(location.Shape)
	if err != nil {
		return err
	}

	_, err = self.conn.Index("location_service", "location", location.Id, nil,
		esLocation{
			ParentId: location.ParentId,
			Name:     location.Name,
		})

	if err != nil {
		return err
	}

	_, err = self.conn.IndexWithParameters("location_service", "geometry", location.Id, location.Id, 0, "", "", "", 0, "", "", false, nil,
		esGeometry{
			Shape: shape,
		})

	return err
}

func (self esStore) FindLocationsByPoint(x, y float64, includeShape bool) ([]model.Location, error) {
	result, err := self.conn.Search("location_service", "location", nil, hash{
		"filter": hash{
			"has_child": hash{
				"type": "geometry",
				"query": hash{
					"geo_shape": hash{
						"shape": hash{
							"shape": hash{
								"type":        "point",
								"coordinates": []float64{x, y},
							},
						},
					},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	locations := make([]model.Location, 0, 5)

	for _, hit := range result.Hits.Hits {
		var location esLocation
		err = json.Unmarshal(*hit.Source, &location)

		if err != nil {
			return nil, err
		}

		locations = append(locations, model.Location{
			Id:       result.Hits.Hits[0].Id,
			ParentId: location.ParentId,
			Name:     location.Name,
		})
	}

	return locations, nil

}

func (self esStore) Begin() Store {
	return self
}

func (self esStore) Flush() {
	self.conn.Flush()
}

func (self esStore) Finish() {
}
