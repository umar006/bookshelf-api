package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	gonanoid "github.com/matoous/go-nanoid/v2"
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

	r.HandleFunc("/books", GetAllBooks).Methods("GET")
	r.HandleFunc("/books/{bookId}", GetBookById).Methods("GET")
	r.HandleFunc("/books", InsertBook).Methods("POST")
	r.HandleFunc("/books/{bookId}", UpdateBookById).Methods("PUT")

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

func GetAllBooks(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	rows, err := db.Queryx("SELECT id, name, publisher FROM book")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var books = []map[string]interface{}{}
	for rows.Next() {
		var book = map[string]interface{}{}
		err = rows.MapScan(book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		books = append(books, book)
	}

	responseData := Response{
		Status: "success",
		Data: struct {
			Books []map[string]interface{} `json:"books"`
		}{books},
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	vars := mux.Vars(r)

	responseData := Response{}

	var book Book
	err := db.QueryRowx("SELECT * FROM book WHERE id=$1", vars["bookId"]).StructScan(&book)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err != nil && err == sql.ErrNoRows {
		responseData.Status = "fail"
		responseData.Message = "Buku tidak ditemukan"

		w.WriteHeader(http.StatusNotFound)
	} else {
		responseData.Status = "success"
		responseData.Data = struct {
			Book Book `json:"book"`
		}{book}
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func UpdateBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	vars := mux.Vars(r)

	var book map[string]any
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := Response{}

	if book["name"] == "" || book["name"] == nil {
		responseData.Status = "fail"
		responseData.Message = "Gagal memperbarui buku. Mohon isi nama buku"

		w.WriteHeader(http.StatusBadRequest)
	} else if book["readPage"].(float64) > book["pageCount"].(float64) {
		responseData.Status = "fail"
		responseData.Message = "Gagal memperbarui buku. readPage tidak boleh lebih besar dari pageCount"

		w.WriteHeader(http.StatusBadRequest)
	} else {
		var bookId string
		err = db.QueryRowx("SELECT id FROM book WHERE id=$1", vars["bookId"]).Scan(&bookId)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if err != nil && err == sql.ErrNoRows {
			responseData.Status = "fail"
			responseData.Message = "Gagal memperbarui buku. Id tidak ditemukan"

			w.WriteHeader(http.StatusNotFound)
		} else {
			book["id"] = bookId

			updateBookQuery := `
                UPDATE book
                SET name=:name,year=:year,author=:author,summary=:summary,publisher=:publisher,
                    page_count=:pageCount,read_page=:readPage,reading=:reading
                WHERE id=:id
            `
			_, err = db.NamedExec(updateBookQuery, book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			responseData.Status = "success"
			responseData.Message = "Buku berhasil diperbarui"
		}
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func DeleteBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	vars := mux.Vars(r)

	_, err := db.Exec("DELETE FROM book WHERE id=$1", vars["bookId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := Response{
		Status:  "success",
		Message: "Buku berhasil dihapus",
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
