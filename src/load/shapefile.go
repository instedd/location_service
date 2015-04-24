package load

import (
	"encoding/csv"
	"geom"
	"github.com/jonas-p/go-shp"
	"io"
	"log"
	"model"
	"os"
	"path/filepath"
	"store"
	"strings"
)

func LoadShapefile(store store.Store, path string, set string, idColumns []string, nameColumn string, defaultTypeName string, typeColumn string, level int) {
	shapefile, err := shp.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer shapefile.Close()

	// Load names from CSV file if exists
	names := make(map[string](string))
	csvPath := strings.Replace(path, filepath.Ext(path), "csv", 1)
	if _, err = os.Stat(csvPath); err == nil {
		namesfile, err := os.Open(csvPath)
		if err != nil {
			log.Fatal(err)
		}

		defer namesfile.Close()
		csvReader := csv.NewReader(namesfile)
		headers, _ := csvReader.Read()

		csvNameIdx := findFieldColumn(headers, nameColumn)
		csvIdColumnsIdx := make([]int, len(idColumns))
		for i, col := range idColumns {
			csvIdColumnsIdx[i] = findFieldColumn(headers, col)
		}

		for err != io.EOF {
			record, err := csvReader.Read()
			if err != nil {
				log.Fatal(err)
				break
			}

			csvIdParts := make([]string, 0, len(idColumns))
			for _, csvIdIdx := range csvIdColumnsIdx {
				csvId := toUtf8(record[csvIdIdx])
				csvIdParts = append(csvIdParts, csvId)
			}
			csvLocationId := set + ":" + strings.Join(csvIdParts, "_")
			csvLocationName := record[csvNameIdx]
			names[csvLocationId] = csvLocationName
		}
	}

	fields := shapefile.Fields()
	fieldNames := make([]string, len(fields))
	for i, field := range fields {
		fieldNames[i] = string(field.Name[:])
	}

	nameIdx := findFieldColumn(fieldNames, nameColumn)
	typeIdx := findFieldColumn(fieldNames, typeColumn)
	idColumnsIdx := make([]int, len(idColumns))
	for i, col := range idColumns {
		idColumnsIdx[i] = findFieldColumn(fieldNames, col)
	}

	for shapefile.Next() {
		n, p := shapefile.Shape()

		// print attributes
		// for k, f := range fields {
		//  val := shapefile.ReadAttribute(n, k)
		//  log.Printf("\t%v: %v\n", f, val)
		// }
		// log.Println()

		idParts := make([]string, 0, len(idColumns))
		for _, idIdx := range idColumnsIdx {
			id := toUtf8(shapefile.ReadAttribute(n, idIdx))
			idParts = append(idParts, id)
		}

		var parentId *string
		if len(idColumns) > 1 {
			str := set + ":" + strings.Join(idParts[:len(idParts)-1], "_")
			parentId = &str
		} else {
			parentId = nil
		}

		locationId := set + ":" + strings.Join(idParts, "_")

		var locationName string
		var found bool
		if locationName, found = names[locationId]; !found {
			locationName = toUtf8(shapefile.ReadAttribute(n, nameIdx))
		}

		typeName := defaultTypeName
		if typeIdx > 0 {
			typeName = toUtf8(shapefile.ReadAttribute(n, typeIdx))
		}

		log.Printf("%s (%s)\n", locationName, typeName)

		shape, err := geom.FromShapefile(p)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		dbLocation := &model.Location{
			Id:       locationId,
			ParentId: parentId,
			Name:     locationName,
			Shape:    model.GeoShape{shape},
			Level:    level,
			TypeName: typeName,
		}

		err = store.AddLocation(dbLocation)
		if err != nil {
			log.Println(err.Error())
			continue
		}

	}

}

func toUtf8(str string) string {
	iso8859_1_buf := []byte(str)
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)
}

func findFieldColumn(fields []string, name string) int {
	if len(name) == 0 {
		return -1
	}

	// Go is always returning false for this comparison, so it never finds the field given its name
	log.Println("Looking for: '" + toUtf8(name) + "'")
	for idx, f := range fields {
		log.Println("Field '" + toUtf8(f) + "' equals?")
		log.Println(toUtf8(f) == toUtf8(name))
		if strings.EqualFold(f, name) {
			log.Println("Found it!")
			return idx
		}
	}

	log.Println(fields)
	log.Fatalf("Field '%s' not found", name)
	return -1
}
