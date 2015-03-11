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

	if location.Id == 0 {
		err := self.db.QueryRow(`INSERT INTO locations(parent_id, level, type_name, name, shape) VALUES ($1, $2, $3, $4, ST_GeomFromWKB($5, 4326)) RETURNING id`,
			location.ParentId, location.Level, location.TypeName, location.Name, wkb).Scan(&location.Id)
		return err
	} else {
		_, err := self.db.Exec(`UPDATE locations SET parent_id = $1, level = $2, type_name = $3, name = $4, shape = ST_GeomFromWKB($5, 4326) WHERE id = $6`,
			location.ParentId, location.Level, location.TypeName, location.Name, wkb, location.Id)
		return err
	}
}

func (self sqlStore) Begin() Store {
	return self
}

func (self sqlStore) Flush() {

}
