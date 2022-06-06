package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Book struct {
	Id        string `json:"bookId"`
	Name      string `json:"name"`
	Year      int16  `json:"year"`
	Author    string `json:"author"`
	Summary   string `json:"summary"`
	Publisher string `json:"pulisher"`
	PageCount int16  `json:"pageCount" db:"page_count"`
	ReadPage  int16  `json:"readPage" db:"read_page"`
	Finished  bool   `json:"finished"`
	Reading   bool   `json:"reading"`
	CreatedAt string `json:"insertedAt"`
	UpdatedAt string `json:"updatedAt"`
}

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

	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
