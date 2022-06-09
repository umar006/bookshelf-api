package book

import (
	"github.com/gorilla/mux"

	"umar006/bookshelf-api/cmd/book/service"
)

func Routes(r *mux.Router) {
	r.HandleFunc("/books", service.GetAllBooks).Methods("GET")
	r.HandleFunc("/books/{bookId}", service.GetBookById).Methods("GET")
	r.HandleFunc("/books", service.InsertBook).Methods("POST")
	r.HandleFunc("/books/{bookId}", service.UpdateBookById).Methods("PUT")
	r.HandleFunc("/books/{bookId}", service.DeleteBookById).Methods("DELETE")
}
