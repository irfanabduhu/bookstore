package review

import (
	"database/sql"

	"github.com/go-chi/chi"
)
func ReviewRouter(db *sql.DB) chi.Router {
	r := chi.NewRouter()
	r.Get("/books/{bookID}/reviews", ListReviewHandler(db))
	r.Post("/books/{bookID}/reviews", CreateReviewHandler(db))
	r.Get("/books/{bookID}/reviews/{reviewID}", GetBookReviewHandler(db))
	r.Put("/books/{bookID}/reviews/{reviewID}", UpdateBookReviewHandler(db))
	r.Delete("/books/{bookID}/reviews/{reviewID}", DeleteBookReviewHandler(db))
	return r
}