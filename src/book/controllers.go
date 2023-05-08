package book

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func ListBooksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, author, price FROM books ORDER BY id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		books := []Book{}
		for rows.Next() {
			book := Book{}
			err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			books = append(books, book)
		}
		if err = rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
	}
}

func CreateBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book Book
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := "INSERT INTO books(title, author, description, price) VALUES ($1, $2, $3, $4) RETURNING id"
		row := db.QueryRow(query, book.Title, book.Author, book.Description, book.Price)
		err = row.Scan(&book.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	}
}

func GetBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book Book

		bookID := chi.URLParam(r, "id")
		query := "SELECT id, title, author, description, price FROM books where id=$1"
		row := db.QueryRow(query, bookID)
		err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.Price)
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	}
}

func UpdateBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the book ID from the URL
		bookID := chi.URLParam(r, "id")

		// Parse the requestBody data from the request body
		var requestBody Book
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var existingBook Book
		query := "SELECT id, title, author, description, price, updated_at FROM books where id=$1"
		row := db.QueryRow(query, bookID)
		err = row.Scan(&existingBook.ID, &existingBook.Title, &existingBook.Author, &existingBook.Description, &existingBook.Price, &existingBook.UpdatedAt)
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if requestBody.Title != "" {
			existingBook.Title = requestBody.Title
		}
		if requestBody.Author != "" {
			existingBook.Author = requestBody.Author
		}
		if requestBody.Description != "" {
			existingBook.Description = requestBody.Description
		}
		if requestBody.Price != nil {
			existingBook.Price = requestBody.Price
		}

		// Update the book in the database
		query = "UPDATE books SET title=$1, author=$2, description=$3, price=$4, updated_at=$5 WHERE id=$6"
		_, err = db.Exec(query, existingBook.Title, existingBook.Author, existingBook.Description, *existingBook.Price, time.Now(), bookID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Return the updated book as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingBook)
	}
}

func DeleteBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "id")

		// Delete the book from the database
		query := "DELETE FROM books WHERE id=$1"
		result, err := db.Exec(query, bookID)
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
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
