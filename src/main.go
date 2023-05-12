package main

import (
	"irfanabduhu/bookstore/book"
	"irfanabduhu/bookstore/config"
	"irfanabduhu/bookstore/review"
	"irfanabduhu/bookstore/user"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	s := config.CreateNewServer()
	s.MountMiddlewares()
	s.MountHandlers(func(r chi.Router) {
		r.Mount("/users", user.UserRouter(s.Databse))
		r.Mount("/books", book.BookRouter(s.Databse))
		r.Mount("/reviews", review.ReviewRouter(s.Databse))
	})
	log.Fatal(http.ListenAndServe(":8080", s.Router))
}