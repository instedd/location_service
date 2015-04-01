package main

import (
	"log"
	"model"
	"store"
	"time"
)

func main() {
	db, err := store.NewSqlStore()
	// db, err := store.NewElasticSearchStore()
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	for i := 0; i < 100; i++ {

		locations, err := db.FindLocationsByPoint(-93.035854, 39.586283, model.ReqOptions{Shapes: false})
		// location, err := db.FindLocationByPoint(-58.575606, -34.608224, false)
		if err != nil {
			log.Fatal(err)
		}

		for _, location := range locations {
			log.Println(location.Name)
		}

		// if (len(locations) > 0) {
		// 	log.Printf("Shape for %s: ", locations[len(locations)-1].Name)
		// 	log.Println(locations[len(locations)-1].Shape)
		// }
	}

	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)

}
