package load

import (
	"geom"
	"github.com/jonas-p/go-shp"
	"log"
	"model"
	"store"
	"strings"
)

func LoadShapefile(store store.Store, path string, set string, idColumns []string, nameColumn string) {
	shapefile, err := shp.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer shapefile.Close()

	fields := shapefile.Fields()
	nameIdx := findFieldColumn(fields, nameColumn)
	idColumnsIdx := make([]int, len(idColumns))
	for i, col := range idColumns {
		idColumnsIdx[i] = findFieldColumn(fields, col)
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
			str := set + ":" + strings.Join(idParts[:len(idParts)-1], "/")
			parentId = &str
		} else {
			parentId = nil
		}

		locationId := set + ":" + strings.Join(idParts, "/")
		locationName := toUtf8(shapefile.ReadAttribute(n, nameIdx))
		log.Println(locationName)

		shape, err := geom.FromShapefile(p)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		dbLocation := &model.Location{
			Id:       locationId,
			ParentId: parentId,
			Name:     locationName,
			Shape:    shape,
		}

		err = store.AddLocation(dbLocation)
		if err != nil {
			log.Fatal(err)
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

func findFieldColumn(fields []shp.Field, name string) int {
	var nameBytes [11]byte
	copy(nameBytes[:], []byte(name))

	for idx, f := range fields {
		if f.Name == nameBytes {
			return idx
		}
	}

	log.Println(fields)
	log.Fatalf("Field '%s' not found", name)
	return -1
}