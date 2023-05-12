package user_test

import (
	"bytes"
	"encoding/json"
	"irfanabduhu/bookstore/utils"
	"net/http"
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
