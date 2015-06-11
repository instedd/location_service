package load

import (
	"bytes"
	"encoding/csv"
	"geom"
	"github.com/jonas-p/go-shp"
	"io"
	"log"
	"model"
	"os"
	"path/filepath"
	"regexp"
	"store"
	"strconv"
	"strings"
)

var csvUnicodePattern, _ = regexp.Compile("<U\\+[0-9A-F]{4}>")

func LoadShapefile(store store.Store, path string, set string, idColumns []string, nameColumn string, defaultTypeName string, typeColumn string, level int) {
	shapefile, err := shp.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer shapefile.Close()

	// Load names from CSV file if exists
	names := make(map[string](string))
	csvPath := strings.Replace(path, filepath.Ext(path), ".csv", 1)
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

		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
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
			csvLocationName = csvUnicodePattern.ReplaceAllStringFunc(csvLocationName, func(str string) string {
				ch, _ := strconv.ParseInt(str[3:7], 16, 0)
				return string(ch)
			})
			names[csvLocationId] = csvLocationName
		}
	}

	fields := shapefile.Fields()
	fieldNames := make([]string, len(fields))
	for i, field := range fields {
		fieldNames[i] = string(field.Name[:bytes.IndexByte(field.Name[:], 0)])
	}

	nameIdx := findFieldColumn(fieldNames, truncate(nameColumn, 10))
	typeIdx := findFieldColumn(fieldNames, truncate(typeColumn, 10))
	idColumnsIdx := make([]int, len(idColumns))
	for i, col := range idColumns {
		idColumnsIdx[i] = findFieldColumn(fieldNames, truncate(col, 10))
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

func truncate(str string, length int) string {
	if len(str) > length {
		return str[:length]
	} else {
		return str
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

	for idx, f := range fields {
		if strings.EqualFold(f, name) {
			return idx
		}
	}

	log.Println(fields)
	log.Fatalf("Field '%s' not found", name)
	return -1
}
