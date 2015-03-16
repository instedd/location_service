package main

import (
	"load"
	"log"
	"store"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [options] FILES:\n", os.Args[0])
		flag.PrintDefaults()
	}

	var source, id, storetype string

	flag.StringVar(&source, "source", "", "Data source (required). Currently supports gadm and naturalearth.")
	flag.StringVar(&id, "id", "", "Identifier prefix for all imported data, defaults to shortened source name.")
	flag.StringVar(&storetype, "store", "sql", "Store type to use. Available: sql, elasticsearch.")

	flag.Parse()

	if (source == "") {
		fmt.Fprintf(os.Stderr, "Error: Option 'source' is required\n")
		flag.Usage()
		return
	}

	var db store.Store
	var err error
	switch storetype {
		case "sql", "psql":
			db, err = store.NewSqlStore()
		case "elastic", "es", "elasticsearch":
			db, err = store.NewElasticSearchStore()
	}

	if err != nil {
		log.Fatal(err)
		return
	}

	for _, file := range flag.Args() {
		var defaultId, nameField string
		var idFields []string
		var level int
		loadSourceSettings(file, source, &defaultId, &idFields, &nameField, &level)
		if (id == "") {
			id = defaultId
		}

		fmt.Printf("\nProcessing %s file %s:\n ID prefix: %s\n ID fields: %s\n Name field: %s\n Level: %d\n\n", source, file, id, idFields, nameField, level)
		load.LoadShapefile(db, file, id, idFields, nameField, level)
	}

	fmt.Printf("\nRunning final update on all imported data\n")
	db.Finish()
}

func loadSourceSettings(file string, source string, defaultId *string, idFields *[]string, nameField *string, plevel *int) {
	switch (source) {

	case "gadm", "GADM":
		*defaultId = "gadm"
		gadmRegexp := regexp.MustCompile("(?i)^(?:.+/)?[a-z]+_adm(\\d+)\\.shp$")
		match := gadmRegexp.FindStringSubmatch(file)
		level, _ := strconv.Atoi(match[1])
		*plevel = level

		if (level == 0) {
			*idFields = []string{"ISO"}
			*nameField = "NAME_ENGLI"
		} else {
			*idFields = make([]string, level+1)
			(*idFields)[0] = "ISO"
			for i := 1; i <= level; i++ {
				(*idFields)[i] = "ID_" + strconv.Itoa(i)
			}
			*nameField = "NAME_" + strconv.Itoa(level)
		}

	case "naturalearth", "NATURALEARTH", "ne", "NE":
		*defaultId = "ne"
		neRegexp := regexp.MustCompile("(?i)^(?:.+/)?[a-z0-9_]+_admin_(\\d+)_[a-z0-9_]+\\.shp$")
		match := neRegexp.FindStringSubmatch(file)
		level, _ := strconv.Atoi(match[1])
		*plevel = level

		if (level == 0) {
			*idFields = []string{"SOV_A3"}
			*nameField = "ADMIN"
		} else if (level == 1) {
			*idFields = []string{"sov_a3", "diss_me"}
			*nameField = "name"
		} else {
			panic("Level " + match[1] + " not supported for natural earth data")
		}

	}
}
