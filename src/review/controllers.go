package review

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
)



func ListBookReviewHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := "SELECT id, user_id, book_id, rating, comment FROM reviews ORDER BY updated_at DESC"
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		reviews := []Review{}
		for rows.Next() {
			review := Review{}
			err := rows.Scan(&review.ID, &review.UserID, &review.BookID, &review.Rating, &review.Comment)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reviews = append(reviews, review)
		}
		if err = rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reviews)
	}
}

func GetBookReviewHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var review Review

		reviewID := chi.URLParam(r, "reviewID")
		query := "SELECT id, user_id, book_id, rating, comment FROM reviews WHERE id = $1"
		row := db.QueryRow(query, reviewID)
		err := row.Scan(&review.ID, &review.UserID, &review.BookID, &review.Rating, &review.Comment)
		if err == sql.ErrNoRows {
			http.Error(w, "Review not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(review)
	}
}


func CreateReviewHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(jwt.MapClaims)
		if !ok {
			http.Error(w, "Missing JWT claims", http.StatusUnauthorized)
			return
		}

		userID, _ := strconv.Atoi(claims["user_id"].(string))

		var review Review
		err := json.NewDecoder(r.Body).Decode(&review)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		review.UserID = userID

		query := "INSERT INTO reviews (user_id, book_id, rating, comment) VALUES ($1, $2, $3, $4) RETURNING id"
		row := db.QueryRow(query, review.UserID, review.BookID, review.Rating, review.Comment)
		err = row.Scan(&review.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(review)
	}
}

func UpdateBookReviewHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(jwt.MapClaims)
		if !ok {
			http.Error(w, "Missing JWT claims", http.StatusUnauthorized)
			return
		}
		userID, _ := strconv.Atoi(claims["user_id"].(string))
		reviewID, _ := strconv.Atoi(chi.URLParam(r, "reviewID"))

		var requestBody Review
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var existingReview Review
		query := "SELECT id, user_id, book_id, rating, comment FROM reviews WHERE id = $1"
		row := db.QueryRow(query, reviewID)
		err = row.Scan(&existingReview.ID, &existingReview.UserID, &existingReview.BookID, &existingReview.Rating, &existingReview.Comment)
		if err == sql.ErrNoRows {
			http.Error(w, "Review not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if userID != existingReview.UserID{
			http.Error(w, "Unauthorized access", http.StatusForbidden)
			return
		}

		if requestBody.Rating != 0 {
			existingReview.Rating = requestBody.Rating
		}
		if requestBody.Comment != "" {
			existingReview.Comment = requestBody.Comment
		}

		query = "UPDATE reviews SET rating=$1, comment=$2 WHERE id=$3"
		_, err = db.Exec(query, existingReview.Rating, existingReview.Comment, reviewID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingReview)
	}
}
func DeleteBookReviewHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(jwt.MapClaims)
		if !ok {
			http.Error(w, "Missing JWT claims", http.StatusUnauthorized)
			return
		}
		userID, _ := strconv.Atoi(claims["user_id"].(string))
		reviewID, _ := strconv.Atoi(chi.URLParam(r, "reviewID"))

		query := "DELETE FROM reviews WHERE id=$1 and user_id=$2"
		result, err := db.Exec(query, reviewID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, "Review not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
