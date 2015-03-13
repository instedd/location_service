package main

import (
	"load"
	"log"
	"store"
)

func main() {
	db, err := store.NewSqlStore()
	// db, err := store.NewElasticSearchStore()
	if err != nil {
		log.Fatal(err)
	}

	load.LoadShapefile(db, "/Users/waj/Downloads/ARG_adm/ARG_adm0.shp", "gadm", []string{"ISO"}, "NAME_ENGLI")
	load.LoadShapefile(db, "/Users/waj/Downloads/ARG_adm/ARG_adm1.shp", "gadm", []string{"ISO", "ID_1"}, "NAME_1")
	load.LoadShapefile(db, "/Users/waj/Downloads/ARG_adm/ARG_adm2.shp", "gadm", []string{"ISO", "ID_1", "ID_2"}, "NAME_2")
}
