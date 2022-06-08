package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func ConnectDB() *sqlx.DB {
	connStr := "postgres://root:root@localhost/bookshelf?sslmode=disable"
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("DB Connected!")

	return db
}
