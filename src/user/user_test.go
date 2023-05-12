package user_test

import (
	"bytes"
	"encoding/json"
	"irfanabduhu/bookstore/utils"
	"net/http"
	"os"
	"testing"
)

func TestAdminLogin(t *testing.T) {
	tests := []struct {
		payload            map[string]string
		expectedStatusCode int
	}{
		{
			map[string]string{
				"username": "admin",
				"password": "abracadabra",
			},
			http.StatusOK,
		},
		{
			map[string]string{
				"username": "admin",
				"password": "anything",
			},
			http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		jsonPayload, _ := json.Marshal(test.payload)
		response := utils.GetResponse(
			"POST",
			"http://localhost:8080/api/v1/users/signin",
			"",
			bytes.NewBuffer(jsonPayload),
		)
		utils.CheckResponseCode(t, test.expectedStatusCode, response.Code)
	}
}

func TestUserSignUp(t *testing.T) {
	user := map[string]string{
		"name": "Irfanul Hoque",
		"username": "irfan",
		"email": "irfan@example.com",
		"password": "123456",
		"plan": "basic",
	}
	jsonPayload, _ := json.Marshal(user)
	response := utils.GetResponse(
		"POST",
		"http://localhost:8080/api/v1/users/signup",
		"",
		bytes.NewBuffer(jsonPayload),
	)
	utils.CheckResponseCode(t, http.StatusCreated, response.Code)
}

func TestUserSignIn(t *testing.T) {
	credentials := map[string]string{
		"username": "irfan",
		"password": "123456",
	}
	jsonPayload, _ := json.Marshal(credentials)
	response := utils.GetResponse(
		"POST",
		"http://localhost:8080/api/v1/users/signin",
		"",
		bytes.NewBuffer(jsonPayload),
	)
	utils.CheckResponseCode(t, http.StatusOK, response.Code)
}

func TestMain(m *testing.M) {
	utils.InitDB()
	code := m.Run()
	utils.TearDown()
    os.Exit(code)
}