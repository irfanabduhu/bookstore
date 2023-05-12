package utils

import (
	"io"
	"irfanabduhu/bookstore/auth"
	"irfanabduhu/bookstore/book"
	"irfanabduhu/bookstore/config"
	"irfanabduhu/bookstore/review"
	"irfanabduhu/bookstore/user"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
)

func InitDB() {
	var queries []string = []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			username VARCHAR(255) UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role VARCHAR(10) DEFAULT 'user' CHECK (role IN ('user', 'admin')),
			plan VARCHAR(10) DEFAULT 'basic' CHECK (plan IN ('basic', 'premium')),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			description TEXT NOT NULL,
			price DECIMAL NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS reviews (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id),
			book_id INTEGER NOT NULL REFERENCES books(id),
			rating INTEGER NOT NULL CHECK (
				rating BETWEEN 1 AND 5
			),
			comment TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`,
		`INSERT INTO users (name, username, email, password, role)
		VALUES (
			'admin',
			'admin',
			'admin@example.com',
			'$2a$10$EOsoyng3jonP9XHiZ3uw5egQAO7Ae0v9Ty75mA0tCU6Z8T9Xf2nj6', -- hash for 'abracadabra' with defaultCost
			'admin'
		);`,
	}

	db := config.ConnectDB()
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Print(err)
		}
	}
}

func TearDown() {
	db := config.ConnectDB()
	queries := []string{
		`DROP TABLE IF EXISTS users CASCADE;`,
		`DROP TABLE IF EXISTS books CASCADE;`,
		`DROP TABLE IF EXISTS reviews CASCADE;`,
	}
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("failed cleaning up db: %v", err.Error())
		} else {
			log.Printf("cleaned up table.")
		}
	}
}

func ExecuteRequest(req *http.Request, s *config.Server) *httptest.ResponseRecorder {
	// create a new ResponseRecorder
	rr := httptest.NewRecorder()
	// execute the request; the handler will write response to the ResponseRecorder
	s.Router.ServeHTTP(rr, req)
	// return it for further inspection
	return rr
}

func GetResponse(method, url string, token string, body io.Reader) *httptest.ResponseRecorder {
	s := config.CreateNewServer()
	s.MountHandlers(func(r chi.Router) {
		r.Mount("/users", user.UserRouter(s.Databse))
		r.Mount("/books", book.BookRouter(s.Databse))
		r.Mount("/reviews", review.ReviewRouter(s.Databse))
	})
	req, _ := http.NewRequest(method, url, body)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	response := ExecuteRequest(req, s)
	return response
}

func CheckResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func GenerateAdminToken() string {
	claims := jwt.MapClaims{
		"user_id":  "1",
		"username": "admin",
		"role":     "admin",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(auth.JWT_SECRET)
	return tokenString
}

func GenerateUserToken() string {
	claims := jwt.MapClaims{
		"user_id":  "2",
		"username": "irfan",
		"role":     "user",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(auth.JWT_SECRET)
	return tokenString
}
