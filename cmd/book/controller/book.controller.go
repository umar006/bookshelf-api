package controller

import (
	"encoding/json"
	"net/http"

	"umar006/bookshelf-api/cmd/book/model"
	"umar006/bookshelf-api/cmd/book/service"
	"umar006/bookshelf-api/pkg"
)

func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var book model.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responseData pkg.Response

	if book.Name == "" {
		responseData.Status = "fail"
		responseData.Message = "Gagal menambahkan buku. Mohon isi nama buku"

		w.WriteHeader(http.StatusBadRequest)
	} else if book.ReadPage > book.PageCount {
		responseData.Status = "fail"
		responseData.Message = "Gagal menambahkan buku. readPage tidak boleh lebih besar dari pageCount"

		w.WriteHeader(http.StatusBadRequest)
	} else {
		_, err = service.InsertBook(&book)
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

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	books, err := service.GetAllBooks(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := pkg.Response{
		Status: "success",
		Data: struct {
			Books []model.Book `json:"books"`
		}{books},
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
