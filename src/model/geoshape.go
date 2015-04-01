package model

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"github.com/foobaz/geom"
	"github.com/foobaz/geom/encoding/wkb"
	_ "github.com/lib/pq"
)

type GeoShape struct {
	geom.T
}

func (g *GeoShape) Scan(src interface{}) error {
	var binshape []byte
	binshape, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}

	shape, _ := wkb.Decode(binshape)
	(*g) = GeoShape{shape}
	return nil
}

func (g *GeoShape) Value() (driver.Value, error) {
	return wkb.Encode((*g).T, binary.LittleEndian, geom.TwoD)
}
