package book_test

import (
	"bytes"
	"encoding/json"
	"irfanabduhu/bookstore/utils"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateNewBook(t *testing.T) {	
	payload := map[string]interface{}{
		"title":       "1984",
		"author":      "George Orwell",
		"description": "To Winston Smith, a young man who works in the Ministry of Truth (Minitru for short), come two people who transform this life completely. One is Julia, whom he meets after she hands him a slip reading, \"I love you.\" The other is O'Brien, who tells him, \"We shall meet in the place where there is no darkness.\" The way in which Winston is betrayed by the one and, against his own desires and instincts, ultimately betrays the other, makes a story of mounting drama and suspense.",
		"price":       2.1,
	}
	jsonPayload, _ := json.Marshal(payload)
	adminToken := utils.GenerateAdminToken()
	userToken := utils.GenerateUserToken()
	response := utils.GetResponse(
		"POST",
		"http://localhost:8080/api/v1/books",
		adminToken,
		bytes.NewBuffer(jsonPayload),
	)
	utils.CheckResponseCode(t, http.StatusCreated, response.Code)

	response = utils.GetResponse(
		"POST",
		"http://localhost:8080/api/v1/books",
		userToken,
		bytes.NewBuffer(jsonPayload),
	)
	utils.CheckResponseCode(t, http.StatusForbidden, response.Code)
}

func TestGetAllBooks(t *testing.T) {
	response := utils.GetResponse(
		"GET",
		"http://localhost:8080/api/v1/books/",
		"",
		nil,
	)
	utils.CheckResponseCode(t, http.StatusOK, response.Code)
}

func TestGetBookDetails(t *testing.T) {
	response := utils.GetResponse(
		"GET",
		"http://localhost:8080/api/v1/books/1",
		"",
		nil,
	)
	utils.CheckResponseCode(t, http.StatusOK, response.Code)
}

type Book struct {
	ID          int       `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Author      string    `json:"author,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       *float64  `json:"price,omitempty"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}


func TestUpdateBook(t *testing.T) {
	payload := map[string]interface{}{
		"title": "mock mock",
	}
	jsonPayload, _ := json.Marshal(payload)
	response := utils.GetResponse(
		"PUT",
		"http://localhost:8080/api/v1/books/1",
		utils.GenerateAdminToken(),
		bytes.NewBuffer(jsonPayload),
	)
	utils.CheckResponseCode(t, http.StatusOK, response.Code)

	var book Book
	json.NewDecoder(response.Body).Decode(&book)
	require.Equal(t, "mock mock", book.Title)

	response = utils.GetResponse(
		"PUT",
		"http://localhost:8080/api/v1/books/1",
		utils.GenerateUserToken(),
		bytes.NewBuffer(jsonPayload),
	)
	utils.CheckResponseCode(t, http.StatusForbidden, response.Code)

	response = utils.GetResponse(
		"PUT",
		"http://localhost:8080/api/v1/books/132",
		utils.GenerateAdminToken(),
		bytes.NewBuffer(jsonPayload),
	)
	utils.CheckResponseCode(t, http.StatusNotFound, response.Code)
}

func TestDeleteBook(t *testing.T) {
	response := utils.GetResponse(
		"DELETE",
		"http://localhost:8080/api/v1/books/1",
		utils.GenerateUserToken(),
		nil,
	)
	utils.CheckResponseCode(t, http.StatusForbidden, response.Code)

	response = utils.GetResponse(
		"DELETE",
		"http://localhost:8080/api/v1/books/1",
		utils.GenerateAdminToken(),
		nil,
	)
	utils.CheckResponseCode(t, http.StatusNoContent, response.Code)
}