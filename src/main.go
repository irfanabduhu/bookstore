package main

import (
	"irfanabduhu/bookstore/book"
	"irfanabduhu/bookstore/config"
	"irfanabduhu/bookstore/review"
	"irfanabduhu/bookstore/user"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	db := config.ConnectDB()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hei buddy! Just chekin'"))
		})
		r.Mount("/users", user.UserRouter(db))
		r.Mount("/books", book.BookRouter(db))
		r.Mount("/reviews", review.ReviewRouter(db))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
