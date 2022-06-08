package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"umar006/bookshelf-api/cmd/book"
)

func main() {
	r := mux.NewRouter()

	book.Routes(r)

	fmt.Println("server started at localhost:9000")
	http.ListenAndServe(":9000", r)
}
