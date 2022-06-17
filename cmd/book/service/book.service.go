package service

import (
	"database/sql"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"umar006/bookshelf-api/cmd/book/model"
	dbx "umar006/bookshelf-api/db"
)

var db *sqlx.DB = dbx.ConnectDB()

func InsertBook(book *model.Book) (int64, error) {
	book.InitId()
	book.InitFinished()

	createBookQuery := `
            INSERT INTO book(id, name, year, author, summary, publisher, page_count, read_page, reading, finished)
            VALUES (:id, :name, :year, :author, :summary, :publisher, :page_count, :read_page, :reading, :finished)
        `
	result, err := db.NamedExec(createBookQuery, &book)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func GetAllBooks(queryParams url.Values) ([]model.Book, error) {
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
		return nil, err
	}

	var books = []model.Book{}
	for result.Next() {
		var book = model.Book{}
		err = result.StructScan(&book)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, err
}

func GetBookById(bookId string) (*model.Book, error) {
	var book model.Book
	err := db.QueryRowx("SELECT * FROM book WHERE id=$1", bookId).StructScan(&book)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}

	return &book, err
}

func UpdateBookById(book *model.Book) (int64, error) {
	book.InitFinished()

	updateBookQuery := `
        UPDATE book
        SET name=:name,year=:year,author=:author,summary=:summary,publisher=:publisher,
            page_count=:page_count,read_page=:read_page,reading=:reading,finished=:finished
        WHERE id=:id
    `
	result, err := db.NamedExec(updateBookQuery, &book)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func DeleteBookById(bookId string) (int64, error) {
	result, err := db.Exec("DELETE FROM book WHERE id=$1", bookId)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
