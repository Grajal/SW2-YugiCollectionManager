package api_tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	api_clients "github.com/Grajal/SW2-YugiCollectionManager/backend/tests/clients"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/tests/factories"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	client := api_clients.NewTestClient(false)

	payload := map[string]interface{}{
		"username": "testuser123",
		"email":    "testuser123@example.com",
		"password": "password123",
	}

	res := client.PerformRequest("POST", "/api/users/", payload, nil)

	assert.Equal(t, http.StatusCreated, res.Code)

	var body map[string]models.User
	err := json.Unmarshal(res.Body.Bytes(), &body)
	assert.NoError(t, err)

	user := body["user"]
	assert.Equal(t, "testuser123", user.Username)
	assert.Equal(t, "testuser123@example.com", user.Email)
	assert.Empty(t, user.Password)
}

func TestGetUserByName_Success(t *testing.T) {
	client := api_clients.NewTestClient(false)
	user := factories.UserFactory()

	res := client.PerformRequest("GET", "/api/users/"+user.Username, nil, nil)

	assert.Equal(t, http.StatusOK, res.Code)

	var body map[string]models.User
	err := json.Unmarshal(res.Body.Bytes(), &body)
	assert.NoError(t, err)

	returnedUser := body["user"]
	assert.Equal(t, user.Username, returnedUser.Username)
	assert.Equal(t, user.Email, returnedUser.Email)
	assert.Empty(t, returnedUser.Password)
}

func TestGetUserByName_NotFound(t *testing.T) {
	client := api_clients.NewTestClient(false)

	res := client.PerformRequest("GET", "/api/users/unknownuser", nil, nil)

	assert.Equal(t, http.StatusNotFound, res.Code)

	var body map[string]string
	err := json.Unmarshal(res.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Contains(t, body["error"], "user not found")
}

func TestDeleteUser_Success(t *testing.T) {
	client := api_clients.NewTestClient(false)
	user := factories.UserFactory()

	res := client.PerformRequest("DELETE", "/api/users/"+user.Username, nil, nil)

	assert.Equal(t, http.StatusOK, res.Code)

	var body map[string]string
	err := json.Unmarshal(res.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "user deleted successfully", body["message"])
}

func TestDeleteUser_NotFound(t *testing.T) {
	client := api_clients.NewTestClient(false)

	res := client.PerformRequest("DELETE", "/api/users/doesnotexist", nil, nil)

	assert.Equal(t, http.StatusNotFound, res.Code)

	var body map[string]string
	err := json.Unmarshal(res.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Contains(t, body["error"], "user not found")
}
