package store

import (
	"database/sql"
	"encoding/binary"
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
	wkb, err := wkb.Encode(location.Shape, binary.LittleEndian, 2)
	if err != nil {
		return err
	}

	_, err = self.db.Exec(`INSERT INTO locations(id, parent_id, level, type_name, name, shape) VALUES ($1, $2, $3, $4, $5, ST_GeomFromWKB($6, 4326))`,
		location.Id, location.ParentId, location.Level, location.TypeName, location.Name, wkb)
	return err
}

func (self sqlStore) Begin() Store {
	return self
}

func (self sqlStore) Flush() {

}
