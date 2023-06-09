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
		r.Use(auth.CurrentUserOnly)
		
		r.Get("/", GetUserHandler(db))
		r.Put("/", UpdateUserHandler(db))
		r.Put("/plan", UpdateUserPlanHandler(db))
		r.Delete("/", DeleteUserHandler(db))
	})

	return r
}