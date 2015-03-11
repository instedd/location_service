package main

import (
	"geom"
	"github.com/jonas-p/go-shp"
	"log"
	"model"
	"store"
)

func main() {
	db, err := store.NewSqlStore()
	// db, err := store.NewElasticSearchStore()
	if err != nil {
		log.Fatal(err)
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
		// 	log.Printf("\t%v: %v\n", f, val)
		// }
		// log.Println()

		locationName := toUtf8([]byte(shapefile.ReadAttribute(n, name_idx)))
		log.Println(locationName)

		shape, err := geom.FromShapefile(p)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		dbLocation := &model.Location{
			Name:  locationName,
			Shape: shape,
		}

		err = db.AddLocation(dbLocation)
		if err != nil {
			log.Fatal(err)
		}

	}

}

func toUtf8(iso8859_1_buf []byte) string {
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)
}
