package geojson

import (
	"encoding/json"
)

type Shape interface {
}

type Point struct {
	X float64
	Y float64
}

type Ring []Point

type Polygon struct {
	Coordinates []Ring
}

type MultiPolygon struct {
	Coordinates [][]Ring
}

func (p Point) MarshalJSON() ([]byte, error) {
	return json.Marshal([]float64{p.X, p.Y})
}

func (p *Polygon) MarshalJSON() ([]byte, error) {
	type polygon struct {
		Type        string `json:"type"`
		Coordinates []Ring `json:"coordinates"`
	}

	return json.Marshal(polygon{
		Type:        "polygon",
		Coordinates: p.Coordinates,
	})
}

func (p *MultiPolygon) MarshalJSON() ([]byte, error) {
	type multipolygon struct {
		Type        string   `json:"type"`
		Coordinates [][]Ring `json:"coordinates"`
	}

	return json.Marshal(multipolygon{
		Type:        "multipolygon",
		Coordinates: p.Coordinates,
	})
}

func (area *Ring) Area() float64 {
	sum := 0.0
	var ax, ay, bx, by, dx, dy float64

	for i, p := range *area {
		if i == 0 {
			ax = 0.0
			ay = 0.0
			dx = -p.Y
			dy = -p.X
		} else {
			ax = p.Y + dx
			ay = p.X + dy
			sum += ax*by - bx*ay
		}
		bx = ax
		by = ay
	}
	return sum / 2
}
