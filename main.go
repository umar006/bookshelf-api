package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	gonanoid "github.com/matoous/go-nanoid/v2"
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

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
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

	responseData := Response{}

	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if book.Name == "" {
		responseData.Status = "fail"
		responseData.Message = "Gagal menambahkan buku. Mohon isi nama buku"

		w.WriteHeader(http.StatusBadRequest)
	} else if book.ReadPage > book.PageCount {
		responseData.Status = "fail"
		responseData.Message = "Gagal menambahkan buku. readPage tidak boleh lebih besar dari pageCount"

		w.WriteHeader(http.StatusBadRequest)
	} else {
		book.Id, err = gonanoid.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		createBookQuery := `
            INSERT INTO book(id, name, year, author, summary, publisher, page_count, read_page, reading)
            VALUES (:id, :name, :year, :author, :summary, :publisher, :page_count, :read_page, :reading)
        `
		_, err = db.NamedExec(createBookQuery, &book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseData.Status = "success"
		responseData.Message = "Buku berhasil ditambahkan"
		responseData.Data = struct {
			BookId string `json:"bookId"`
		}{book.Id}

		w.WriteHeader(http.StatusCreated)
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
