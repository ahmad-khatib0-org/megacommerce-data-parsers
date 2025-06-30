package main

import (
	"context"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatal(err)
	}

	db, err := OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = InsertCategories(ctx, tx)

	if err == nil {
		err = tx.Commit(ctx)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println(err)
		err := tx.Rollback(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
}
