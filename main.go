package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func main() {
	var err error
	connStr := "postgres://root:root@localhost/bookshelf?sslmode=disable"
	db, err = sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("DB Connected!")

	r := mux.NewRouter()

	r.HandleFunc("/books", InsertBook).Methods("POST")

	fmt.Println("server started at localhost:9000")
	http.ListenAndServe(":9000", r)
}

func InsertBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
}
