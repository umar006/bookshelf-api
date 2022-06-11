package service

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"umar006/bookshelf-api/cmd/book/model"
	dbx "umar006/bookshelf-api/db"
	"umar006/bookshelf-api/pkg"
)

var db *sqlx.DB = dbx.ConnectDB()

func InsertBook(book *model.Book) (sql.Result, error) {
	book.Id, _ = gonanoid.New()
	book.Finished = book.PageCount == book.ReadPage

	createBookQuery := `
            INSERT INTO book(id, name, year, author, summary, publisher, page_count, read_page, reading, finished)
            VALUES (:id, :name, :year, :author, :summary, :publisher, :page_count, :read_page, :reading, :finished)
        `
	result, err := db.NamedExec(createBookQuery, &book)

	return result, err
}

func GetAllBooks(queryParams url.Values) (*sqlx.Rows, error) {
	var result *sqlx.Rows
	var err error
	getBooksQuery := `
        SELECT id, name, publisher
        FROM book
    `
	if queryParams.Has("reading") {
		getBooksQuery += "WHERE reading=$1"
		result, err = db.Queryx(getBooksQuery, queryParams.Get("reading"))
	} else if queryParams.Has("finished") {
		getBooksQuery += "WHERE finished=$1"
		result, err = db.Queryx(getBooksQuery, queryParams.Get("finished"))
	} else if queryParams.Has("name") {
		getBooksQuery += "WHERE name ILIKE $1"
		result, err = db.Queryx(getBooksQuery, "%"+queryParams.Get("name")+"%")
	} else {
		result, err = db.Queryx(getBooksQuery)
	}

	if err != nil {
		return result, err
	}

	return result, err
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	vars := mux.Vars(r)

	var responseData pkg.Response

	var book model.Book
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
			Book model.Book `json:"book"`
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

	var responseData pkg.Response

	if book["name"] == "" || book["name"] == nil {
		responseData.Status = "fail"
		responseData.Message = "Gagal memperbarui buku. Mohon isi nama buku"

		w.WriteHeader(http.StatusBadRequest)
	} else if book["readPage"].(float64) > book["pageCount"].(float64) {
		responseData.Status = "fail"
		responseData.Message = "Gagal memperbarui buku. readPage tidak boleh lebih besar dari pageCount"

		w.WriteHeader(http.StatusBadRequest)
	} else {
		book["id"] = vars["bookId"]

		updateBookQuery := `
            UPDATE book
            SET name=:name,year=:year,author=:author,summary=:summary,publisher=:publisher,
                page_count=:pageCount,read_page=:readPage,reading=:reading
            WHERE id=:id
        `
		result, err := db.NamedExec(updateBookQuery, book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		affectedRow, _ := result.RowsAffected()
		if affectedRow == 1 {
			responseData.Status = "success"
			responseData.Message = "Buku berhasil diperbarui"
		} else {
			responseData.Status = "fail"
			responseData.Message = "Gagal memperbarui buku. Id tidak ditemukan"

			w.WriteHeader(http.StatusNotFound)
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

	result, err := db.Exec("DELETE FROM book WHERE id=$1", vars["bookId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responseData pkg.Response

	affectedRow, _ := result.RowsAffected()
	if affectedRow == 1 {
		responseData.Status = "success"
		responseData.Message = "Buku berhasil dihapus"
	} else {
		responseData.Status = "fail"
		responseData.Message = "Buku gagal dihapus. Id tidak ditemukan"

		w.WriteHeader(http.StatusNotFound)
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
