package store

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"github.com/foobaz/geom"
	"github.com/foobaz/geom/encoding/wkb"
	_ "github.com/lib/pq"
	"model"
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

func (self sqlStore) AddLocation(location *model.Location) error {
	wkb, err := wkb.Encode(location.Shape, binary.LittleEndian, geom.TwoD)
	if err != nil {
		return err
	}

	var ancestors StringSlice

	if location.ParentId != nil {
		err = self.db.QueryRow("SELECT ancestors_ids FROM locations WHERE id = $1", *location.ParentId).Scan(&ancestors)
		if err != nil {
			return err
		}

		ancestors = append(ancestors, *location.ParentId)
	} else {
		ancestors = make([]string, 0)
	}

	_, err = self.db.Exec(`INSERT INTO locations(id, parent_id, level, type_name, name, shape, ancestors_ids) VALUES ($1, $2, $3, $4, $5, ST_GeomFromWKB($6, 4326), $7)`,
		location.Id, location.ParentId, location.Level, location.TypeName, location.Name, wkb, &ancestors)
	return err
}

func (self sqlStore) FindLocationsByPoint(x, y float64, includeShape bool) ([]model.Location, error) {
	var fields string
	if includeShape {
		fields = `a.id, a.parent_id, a.name, ST_AsBinary(a.shape) as binshape`
	} else {
		fields = `a.id, a.parent_id, a.name`
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM locations AS a
			INNER JOIN locations as l
			ON a.id = ANY(l.ancestors_ids) OR a.id = l.id
		WHERE l.leaf = TRUE
			AND ST_Within(ST_SetSRID(ST_Point($1, $2), 4326), l.shape)
		ORDER BY a.level DESC`, fields)

	rows, err := self.db.Query(query, x, y)
	if err != nil {
		return nil, err
	}

	locations := make([]model.Location, 0, 20)

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
