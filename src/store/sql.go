package store

import (
	"database/sql"
	"encoding/binary"
	"github.com/foobaz/geom/encoding/wkb"
	_ "github.com/lib/pq"
	"model"
	"fmt"
	"github.com/foobaz/geom"
)

type sqlStore struct {
	db *sql.DB
}

func NewSqlStore() (Store, error) {
	db, err := sql.Open("postgres", "dbname=location_service sslmode=disable")
	if err != nil {
		return nil, err
	}

	return sqlStore{db: db}, nil
}

// TODO: Implement Valuer and Scanner interfaces for geometries
// See http://jmoiron.net/blog/built-in-interfaces/
// func (l geom.T) Value() (driver.Value, error) {
// 	wkb, err := wkb.Encode(location.Shape, binary.LittleEndian, 2)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return wkb, nil
// }


func (self sqlStore) AddLocation(location *model.Location) error {
	wkb, err := wkb.Encode(location.Shape, binary.LittleEndian, geom.TwoD)
	if err != nil {
		return err
	}

	_, err = self.db.Exec(`INSERT INTO locations(id, parent_id, level, type_name, name, shape) VALUES ($1, $2, $3, $4, $5, ST_GeomFromWKB($6, 4326))`,
		location.Id, location.ParentId, location.Level, location.TypeName, location.Name, wkb)
	return err
}

func (self sqlStore) FindLocationsByPoint(x, y float64, includeShape bool) ([]model.Location, error) {
	var fields string
	if (includeShape) {
		fields = `id, parent_id, name, ST_AsBinary(shape) as binshape`
	} else {
		fields = `id, parent_id, name`
	}

	query := fmt.Sprintf("SELECT %s FROM locations WHERE ST_Within(ST_SetSRID(ST_Point($1, $2), 4326), shape)", fields)
	rows, err := self.db.Query(query, x, y)
	if err != nil {
		return nil, err
	}

	locations := make([]model.Location, 0, 5)

	for rows.Next() {
		var location model.Location

		var err error
		if includeShape {
			var binshape []byte
			err = rows.Scan(&location.Id, &location.ParentId, &location.Name, &binshape)
			shape, _ := wkb.Decode(binshape)
			location.Shape = shape
		} else {
			err = rows.Scan(&location.Id, &location.ParentId, &location.Name)
		}

		if err != nil {
			return nil, err
		}

		locations = append(locations, location)
	}

	return locations, nil
}

func (self sqlStore) Begin() Store {
	return self
}

func (self sqlStore) Flush() {
}

func (self sqlStore) Finish() {
	self.db.Exec("UPDATE locations SET leaf = NOT EXISTS (SELECT 1 FROM locations l2 WHERE l2.parent_id = locations.id)")
}
