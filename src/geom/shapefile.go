package geom

import (
	"fmt"
	"github.com/foobaz/geom"
	"github.com/jonas-p/go-shp"
	"reflect"
)

func FromShapefile(shape shp.Shape) (geom.T, error) {
	switch v := shape.(type) {
	case *shp.Polygon:
		return FromShapefilePolygon(v), nil
	default:
		return nil, fmt.Errorf("Unsupported shapefile geometry type: %q", reflect.TypeOf(shape))
	}
}

func FromShapefilePolygon(pol *shp.Polygon) geom.T {
	rings := make([]geom.Ring, pol.NumParts)
	for part, part_start := range pol.Parts {
		var part_end int32
		if int32(part) == pol.NumParts-1 {
			part_end = pol.NumPoints - 1
		} else {
			part_end = pol.Parts[part+1] - 1
		}

		rings[part] = make([]geom.Point, part_end-part_start+1)
		for i := part_start; i <= part_end; i++ {
			rings[part][i-part_start] = geom.Point{pol.Points[i].X, pol.Points[i].Y}
		}
	}

	outer := make([]geom.Polygon, 0, pol.NumParts)
	inner := make([]geom.Ring, 0, pol.NumParts)

	for _, ring := range rings {
		if ringArea(&ring) > 0 {
			outer = append(outer, geom.Polygon{ring})
		} else {
			inner = append(inner, ring)
		}
	}

	//TODO: use inner rings

	if len(outer) > 1 {
		return geom.MultiPolygon(outer)
	} else {
		return outer[0]
	}
}

func ringArea(ring *geom.Ring) float64 {
	sum := 0.0
	var ax, ay, bx, by, dx, dy float64

	for i, p := range *ring {
		if i == 0 {
			ax = 0.0
			ay = 0.0
			dx = -p[1]
			dy = -p[0]
		} else {
			ax = p[1] + dx
			ay = p[0] + dy
			sum += ax*by - bx*ay
		}
		bx = ax
		by = ay
	}
	return -sum / 2
}
