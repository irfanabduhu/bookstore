package book

import (
	"database/sql"

	"github.com/go-chi/chi"
)

func BookRouter(db *sql.DB) chi.Router {
	r := chi.NewRouter()
	r.Get("/books", ListBooksHandler(db))
	r.Post("/books", CreateBookHandler(db))
	r.Get("/books/{id}", GetBookHandler(db))
	r.Put("/books/{id}", UpdateBookHandler(db))
	r.Delete("/books/{id}", DeleteBookHandler(db))
	return r
}
