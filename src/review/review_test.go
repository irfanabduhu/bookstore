package review_test

import (
	"bytes"
	"encoding/json"
	"irfanabduhu/bookstore/utils"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateNewReview(t *testing.T) {
	payload := map[string]interface{}{
		"book_id": 3,
		"rating":  4,
		"comment": "Awesome book",
	}
	jsonPayload, _ := json.Marshal(payload)
	response := utils.GetResponse(
		"POST",
		"http://localhost:8080/api/v1/reviews",
		utils.GenerateUserToken(),
		bytes.NewBuffer(jsonPayload),
	)
	utils.CheckResponseCode(t, http.StatusCreated, response.Code)

	response = utils.GetResponse(
		"POST",
		"http://localhost:8080/api/v1/reviews",
		"",
		bytes.NewBuffer(jsonPayload),
	)
	utils.CheckResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestGetAllReviewsByBook(t *testing.T) {
	response := utils.GetResponse(
		"GET",
		"http://localhost:8080/api/v1/reviews/books/1",
		"",
		nil,
	)
	utils.CheckResponseCode(t, http.StatusOK, response.Code)
}

func TestGetReviewDetails(t *testing.T) {
	response := utils.GetResponse(
		"GET",
		"http://localhost:8080/api/v1/reviews/1",
		"",
		nil,
	)
	utils.CheckResponseCode(t, http.StatusOK, response.Code)
}

type Review struct {
	ID        int       `json:"id,omitempty"`
	UserID    int       `json:"user_id,omitempty"`
	BookID    int       `json:"book_id,omitempty"`
	Rating    int       `json:"rating,omitempty"`
	Comment   string    `json:"comment,omitempty"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func TestUpdateReview(t *testing.T) {
	payload := map[string]interface{}{
		"rating": 2,
	}
	jsonPayload, _ := json.Marshal(payload)
	response := utils.GetResponse(
		"PUT",
		"http://localhost:8080/api/v1/reviews/1",
		utils.GenerateUserToken(),
		bytes.NewBuffer(jsonPayload),
	)
	utils.CheckResponseCode(t, http.StatusOK, response.Code)

	var review Review
	json.NewDecoder(response.Body).Decode(&review)
	require.Equal(t, 2, review.Rating)

	response = utils.GetResponse(
		"PUT",
		"http://localhost:8080/api/v1/reviews/1",
		"",
		bytes.NewBuffer(jsonPayload),
	)
	utils.CheckResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestDeleteReview(t *testing.T) {
	response := utils.GetResponse(
		"DELETE",
		"http://localhost:8080/api/v1/reviews/1",
		"",
		nil,
	)
	utils.CheckResponseCode(t, http.StatusUnauthorized, response.Code)

	response = utils.GetResponse(
		"DELETE",
		"http://localhost:8080/api/v1/reviews/1",
		utils.GenerateUserToken(),
		nil,
	)
	utils.CheckResponseCode(t, http.StatusNoContent, response.Code)
}

func initReview() {
	user := map[string]string{
		"name":     "Irfanul Hoque",
		"username": "irfan",
		"email":    "irfan@example.com",
		"password": "123456",
		"plan":     "basic",
	}
	jsonPayload, _ := json.Marshal(user)
	utils.GetResponse(
		"POST",
		"http://localhost:8080/api/v1/users/signup",
		"",
		bytes.NewBuffer(jsonPayload),
	)

	books := []map[string]interface{}{
		{
			"title":       "1984",
			"author":      "George Orwell",
			"description": "To Winston Smith, a young man who works in the Ministry of Truth (Minitru for short), come two people who transform this life completely. One is Julia, whom he meets after she hands him a slip reading, \"I love you.\" The other is O'Brien, who tells him, \"We shall meet in the place where there is no darkness.\" The way in which Winston is betrayed by the one and, against his own desires and instincts, ultimately betrays the other, makes a story of mounting drama and suspense.",
			"price":       2.1,
		},
		{
			"title":       "To Kill a Mockingbird",
			"author":      "Harper Lee",
			"description": "Compassionate, dramatic, and deeply moving.",
			"price":       1,
		},
		{
			"title":       "The Adventures of Huckleberry Finn",
			"author":      "Mark Twain",
			"description": "A nineteenth-century boy from a Mississippi River town recounts his adventures as he travels down the river with a runaway slave, encountering a family involved in a feud, two scoundrels pretending to be royalty, and Tom Sawyer's aunt who mistakes him for Tom.",
			"price":       0,
		},
		{
			"title":       "Les Misérables",
			"author":      "Victor Hugo",
			"description": "Victor Hugo's tale of injustice, heroism and love follows the fortunes of Jean Valjean, an escaped convict determined to put his criminal past behind him.",
			"price":       3.1,
		},
	}
	for _, book := range books {
		jsonPayload, _ := json.Marshal(book)
		utils.GetResponse(
			"POST",
			"http://localhost:8080/api/v1/books",
			utils.GenerateUserToken(),
			bytes.NewBuffer(jsonPayload),
		)
		utils.GetResponse(
			"POST",
			"http://localhost:8080/api/v1/books",
			utils.GenerateAdminToken(),
			bytes.NewBuffer(jsonPayload),
		)
	}
}

func TestMain(m *testing.M) {
	utils.InitDB()
	initReview()
	code := m.Run()
	utils.TearDown()
	os.Exit(code)
}
