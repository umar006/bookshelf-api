package book

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	dbx "umar006/bookshelf-api/db"
)

var db *sqlx.DB

func Routes(r *mux.Router) {
	db = dbx.ConnectDB()

	r.HandleFunc("/books", GetAllBooks).Methods("GET")
	r.HandleFunc("/books/{bookId}", GetBookById).Methods("GET")
	r.HandleFunc("/books", InsertBook).Methods("POST")
	r.HandleFunc("/books/{bookId}", UpdateBookById).Methods("PUT")
	r.HandleFunc("/books/{bookId}", DeleteBookById).Methods("DELETE")
}
