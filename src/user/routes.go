package user

import (
	"database/sql"

	"github.com/go-chi/chi"
)

func UserRouter(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	r.Post("/users/signup", SignUpHandler(db))
	r.Post("/users/signin", SignInHandler(db))
	r.Get("/users/{id}", GetUserHandler(db))
	r.Put("/users/{id}", UpdateUserHandler(db))
	r.Put("/users/{id}/plan", UpdateUserPlanHandler(db))
	r.Post("/users/{id}/rent", RentBookHandler(db))
	r.Post("/users/{id}/buy", BuyBookHandler(db))
	
	return r
}
