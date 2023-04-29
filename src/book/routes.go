package book

import (
	"database/sql"
	"irfanabduhu/bookstore/auth"

	"github.com/go-chi/chi"
)

func BookRouter(db *sql.DB) chi.Router {
	r := chi.NewRouter()
	r.Get("/", ListBooksHandler(db))
	r.Get("/{id}", GetBookHandler(db))

	r.Group(func(r chi.Router) {
		r.Use(auth.Auth)
		r.Use(auth.AdminOnly)

		r.Post("/", CreateBookHandler(db))
		r.Put("/{id}", UpdateBookHandler(db))
		r.Delete("/{id}", DeleteBookHandler(db))
	})

	return r
}
