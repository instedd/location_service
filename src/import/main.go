package main

import (
	"fmt"
	"geojson"
	"github.com/jonas-p/go-shp"
	elastigo "github.com/mattbaird/elastigo/lib"
)

type Location struct {
	Name  string        `json:"name"`
	Shape geojson.Shape `json:"shape"`
}

func main() {
	es := elastigo.NewConn()
	es.CreateIndex("location_service")

	type hash map[string]interface{}
	err := es.PutMapping("location_service", "location", Location{}, elastigo.MappingOptions{
		Properties: hash{"shape": hash{
			"type": "geo_shape",
		}},
	})
	if err != nil {
		println(err.Error())
	}

	shapefile, _ := shp.Open("/Users/waj/Downloads/ARG_adm/ARG_adm1.shp")
	defer shapefile.Close()

	fields := shapefile.Fields()

	var name [11]byte
	copy(name[:], []byte("NAME_1"))

	name_idx := -1
	for idx, f := range fields {
		if f.Name == name {
			name_idx = idx
			break
		}
	}

	for shapefile.Next() {
		n, p := shapefile.Shape()

		// print attributes
		// for k, f := range fields {
		// 	val := shapefile.ReadAttribute(n, k)
		// 	fmt.Printf("\t%v: %v\n", f, val)
		// }
		// fmt.Println()

		shape, err := geojson.FromShapefile(p)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		location := Location{
			Name:  shapefile.ReadAttribute(n, name_idx),
			Shape: shape,
		}

		_, err = es.Index("location_service", "location", "", nil, location)
		if err != nil {
			println(location.Name)
			println(err.Error())
		}

	}

	es.Flush()
}
