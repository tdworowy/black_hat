package main

import (
	"database/sql"
	"fmt"
	"log"
)

// import correct driver
func main() {

	db, err := sql.Open("Postgres", "root:password@tcp(127.0.0.1:3306)/store")
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	var (
		ccnum, date, cvv, exp string
		amount                float64
	)
	rows, err := db.Query("select ccnum, date, cvv, exp, amount from transactions")
	if err != nil {
		log.Panicln(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&ccnum, &date, &cvv, &exp, &amount)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println(ccnum, date, cvv, exp, amount)
	}
	if rows.Err() != nil {
		log.Panicln(err)
	}
}
