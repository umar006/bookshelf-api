package book

import (
	"github.com/gorilla/mux"

	"umar006/bookshelf-api/cmd/book/controller"
)

func Routes(r *mux.Router) {
	r.HandleFunc("/books", controller.GetAllBooks).Methods("GET")
	r.HandleFunc("/books/{bookId}", controller.GetBookById).Methods("GET")
	r.HandleFunc("/books", controller.CreateBook).Methods("POST")
	r.HandleFunc("/books/{bookId}", controller.UpdateBookById).Methods("PUT")
	r.HandleFunc("/books/{bookId}", controller.DeleteBookById).Methods("DELETE")
}
