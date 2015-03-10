package geojson

import (
	"fmt"
	"github.com/jonas-p/go-shp"
	"reflect"
)

func FromShapefile(shape shp.Shape) (Shape, error) {
	switch v := shape.(type) {
	case *shp.Polygon:
		return FromShapefilePolygon(v), nil
	default:
		return nil, fmt.Errorf("Unsupported shapefile geometry type: %q", reflect.TypeOf(shape))
	}
}

func FromShapefilePolygon(pol *shp.Polygon) Shape {
	rings := make([]Ring, pol.NumParts)
	for part, part_start := range pol.Parts {
		var part_end int32
		if int32(part) == pol.NumParts-1 {
			part_end = pol.NumPoints - 1
		} else {
			part_end = pol.Parts[part+1] - 1
		}

		rings[part] = make(Ring, part_end-part_start+1)
		for i := part_start; i <= part_end; i++ {
			rings[part][part_end-i] = Point{pol.Points[i].X, pol.Points[i].Y}
		}
	}

	outer := make([][]Ring, 0, pol.NumParts)
	inner := make([]Ring, 0, pol.NumParts)

	for _, ring := range rings {
		if ring.Area() > 0 {
			outer = append(outer, []Ring{ring})
		} else {
			inner = append(inner, ring)
		}
	}

	//TODO: use inner rings

	if len(outer) > 1 {
		return &MultiPolygon{Coordinates: outer}
	} else {
		return &Polygon{Coordinates: outer[0]}
	}

}
