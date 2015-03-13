package main

import (
	"log"
	"store"
	"time"
)

func main() {
	// db, err := store.NewSqlStore()
	db, err := store.NewElasticSearchStore()
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	for i := 0; i < 100; i++ {

		locations, err := db.FindLocationsByPoint(-93.035854, 39.586283, false)
		// location, err := db.FindLocationByPoint(-58.575606, -34.608224, false)
		if err != nil {
			log.Fatal(err)
		}

		for _, location := range locations {
			log.Println(location.Name)
		}
	}

	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)

}
