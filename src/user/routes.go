package user

import (
	"database/sql"
	"irfanabduhu/bookstore/auth"

	"github.com/go-chi/chi"
)

func UserRouter(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", SignUpHandler(db))
	r.Post("/signin", SignInHandler(db))
	r.Route("/{username}", func(r chi.Router) {
		r.Use(auth.Auth)
		
		r.Get("/", GetUserHandler(db))
		r.Put("/", UpdateUserHandler(db))
		r.Put("/plan", UpdateUserPlanHandler(db))
		r.Post("/rent", RentBookHandler(db))
		r.Post("/buy", BuyBookHandler(db))
	})

	return r
}