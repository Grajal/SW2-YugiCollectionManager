package api_tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	api_clients "github.com/Grajal/SW2-YugiCollectionManager/backend/tests/clients"
	"github.com/stretchr/testify/assert"
)

func TestAuthFlow(t *testing.T) {
	client := api_clients.NewTestClient(false)

	// Test data
	user := map[string]string{
		"username": "testuser",
		"email":    "testuser@example.com",
		"password": "securepassword",
	}

	t.Run("Register new user", func(t *testing.T) {
		resp := client.PerformRequest("POST", "/api/auth/register", user, nil)
		assert.Equal(t, http.StatusCreated, resp.Code)

		var body map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)
		assert.Equal(t, "User registered successfully", body["message"])
	})

	t.Run("Register duplicate user", func(t *testing.T) {
		resp := client.PerformRequest("POST", "/api/auth/register", user, nil)
		assert.Equal(t, http.StatusConflict, resp.Code)

		var body map[string]string
		_ = json.Unmarshal(resp.Body.Bytes(), &body)
		assert.Contains(t, body["error"], "already exists")
	})

	var token string
	t.Run("Login with correct credentials", func(t *testing.T) {
		loginReq := map[string]string{
			"username": user["username"],
			"password": user["password"],
		}
		resp := client.PerformRequest("POST", "/api/auth/login", loginReq, nil)
		assert.Equal(t, http.StatusOK, resp.Code)

		cookies := resp.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "token" {
				token = cookie.Value
			}
		}
		assert.NotEmpty(t, token)
	})

	t.Run("Login with wrong password", func(t *testing.T) {
		loginReq := map[string]string{
			"username": user["username"],
			"password": "wrongpassword",
		}
		resp := client.PerformRequest("POST", "/api/auth/login", loginReq, nil)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)

		var body map[string]string
		_ = json.Unmarshal(resp.Body.Bytes(), &body)
		assert.Contains(t, body["error"], "invalid")
	})

	t.Run("Fetch current user with token", func(t *testing.T) {
		headers := map[string]string{
			"Cookie": "token=" + token,
		}
		resp := client.PerformRequest("GET", "/api/auth/me", nil, headers)
		assert.Equal(t, http.StatusOK, resp.Code)

		var currentUser models.User
		err := json.Unmarshal(resp.Body.Bytes(), &currentUser)
		assert.NoError(t, err)
		assert.Equal(t, user["username"], currentUser.Username)
		assert.Equal(t, user["email"], currentUser.Email)
	})
}
