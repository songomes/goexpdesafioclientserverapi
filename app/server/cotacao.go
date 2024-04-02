package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./cotacao.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, bid FROM cotacao")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var bid string
		err = rows.Scan(&id, &bid)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ID: %d, BID: %s\n", id, bid)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
