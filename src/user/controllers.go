package user

import (
	"database/sql"
	"encoding/json"
	"irfanabduhu/bookstore/auth"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUpHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Hash the password before storing it in the database
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Store the user in the database
		query := "INSERT INTO users (name, username, email, password, plan) VALUES ($1, $2, $3, $4, $5) RETURNING id, role"
		row := db.QueryRow(query, user.Name, user.Username, user.Email, string(hashedPassword), user.Plan)
		err = row.Scan(&user.ID, &user.Role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Hide password from response:
		user.Password = ""
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func SignInHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Retrieve the user from the database
		var user User
		query := "SELECT id, username, password, role FROM users WHERE username=$1"
		row := db.QueryRow(query, creds.Username)
		err = row.Scan(&user.ID, &user.Username, &user.Password, &user.Role)

		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Compare the hashed password with the provided password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Create a JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username":   user.Username,
			"role": user.Role,
			"exp":  time.Now().Add(time.Hour * 24).Unix(),
		})

		// Sign the token with a secret key
		tokenString, err := token.SignedString([]byte(auth.JWT_SECRET))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the token as a JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TokenResponse{Token: tokenString})
	}
}

func GetUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		user.Username = chi.URLParam(r, "username")

		query := "SELECT id, name, email, role, plan FROM users WHERE username=$1"
		row := db.QueryRow(query, user.Username)
		err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Plan)
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func UpdateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse username from URL path
		username := chi.URLParam(r, "username")

		// Parse updated user details from request body
		var requestBody User
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// fetch existing user from the db:
		var existingUser User
		query := "SELECT id, name, username, email, password, role, plan FROM users WHERE username=$1"
		row := db.QueryRow(query, username)
		err = row.Scan(
			&existingUser.ID,
			&existingUser.Name,
			&existingUser.Username,
			&existingUser.Email,
			&existingUser.Password,
			&existingUser.Role,
			&existingUser.Plan,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Update the user's details
		if requestBody.Name != "" {
			existingUser.Name = requestBody.Name
		}
		if requestBody.Username != "" {
			existingUser.Username = requestBody.Username
		}
		if requestBody.Email != "" {
			existingUser.Email = requestBody.Email
		}
		if requestBody.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			existingUser.Password = string(hashedPassword)
		}
		if requestBody.Plan != "" {
			existingUser.Plan = requestBody.Plan
		}

		// Update the user in the database
		query = "UPDATE users SET name=$1, username=$2, email=$3, password=$4, plan=$5, updated_at=$6 WHERE username=$7"
		_, err = db.Exec(
			query,
			existingUser.Name,
			existingUser.Username,
			existingUser.Email,
			existingUser.Password,
			existingUser.Plan,
			time.Now(),
			username,
		)
		if err != nil {
			e := err.(*pq.Error)
			if e.Code == "23505" {
				http.Error(w, "'username' and 'email' must be unique.", http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the updated user as a JSON response
		existingUser.Password = "" // hide the password field
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingUser)
	}
}

func UpdateUserPlanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request:
		username := chi.URLParam(r, "username")
		plan := struct {
			Plan string `json:"plan"`
		}{}
		err := json.NewDecoder(r.Body).Decode(&plan)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Update database:
		query := "UPDATE users SET plan=$1 WHERE username=$2"
		res, err := db.Exec(query, plan.Plan, username)
		if err != nil {
			e := err.(*pq.Error)
			if e.Code == "23514" {
				http.Error(w, "Key 'plan' should have value of either 'basic' or 'premium'.", http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if n, _ := res.RowsAffected(); n == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		query := "DELETE FROM users WHERE username=$1"
		result, err := db.Exec(query, username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if a row was affected by the delete operation
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func RentBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func BuyBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
