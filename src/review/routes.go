package review

import (
	"database/sql"
	"irfanabduhu/bookstore/auth"

	"github.com/go-chi/chi"
)
func ReviewRouter(db *sql.DB) chi.Router {
	r := chi.NewRouter()
	r.Get("/books/{bookID}", ListBookReviewHandler(db))
	r.Get("/{reviewID}", GetBookReviewHandler(db))

	r.Group(func(r chi.Router) {
		r.Use(auth.Auth)

		r.Post("/", CreateReviewHandler(db))
		r.Put("/{reviewID}", UpdateBookReviewHandler(db))
		r.Delete("/{reviewID}", DeleteBookReviewHandler(db))
	})
	return r
}