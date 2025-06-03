package api_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	api_clients "github.com/Grajal/SW2-YugiCollectionManager/backend/tests/clients"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/tests/factories"
	"github.com/stretchr/testify/assert"
)

func TestGetCollection_Success(t *testing.T) {
	client := api_clients.NewTestClient(true) // `true` means user authenticated and user_id set in context

	response := client.PerformRequest("GET", "/api/collections/", nil, nil)

	assert.Equal(t, http.StatusOK, response.Code)

	var body map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Contains(t, body, "collection")
}

func TestGetCollectionCard_Success(t *testing.T) {
	client := api_clients.NewTestClient(true)
	userID := client.User.ID

	card := factories.CardFactory("Blue-Eyes White Dragon")

	factories.UserCardFactory(userID, card.ID, 1)

	response := client.PerformRequest("GET", fmt.Sprintf("/api/collections/%d", card.ID), nil, nil)

	assert.Equal(t, http.StatusOK, response.Code)

	var userCard map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &userCard)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), userCard["Quantity"])
}

func TestGetCollectionCard_InvalidID(t *testing.T) {
	client := api_clients.NewTestClient(true)

	response := client.PerformRequest("GET", "/api/collections/abc", nil, nil)

	assert.Equal(t, http.StatusBadRequest, response.Code)

	var body map[string]string
	json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, "Invalid card ID", body["error"])
}

func TestAddCardToCollection_Success(t *testing.T) {
	client := api_clients.NewTestClient(true)

	card := factories.CardFactory("Blue-Eyes White Dragon")

	payload := map[string]interface{}{
		"card_id":  card.ID,
		"quantity": 2,
	}

	response := client.PerformRequest("POST", "/api/collections/", payload, nil)

	assert.Equal(t, http.StatusOK, response.Code)

	var res map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "Card added to collection successfully", res["message"])
}

func TestAddCardToCollection_InvalidInput(t *testing.T) {
	client := api_clients.NewTestClient(true)

	jsonBody := `{"card_id":0,"quantity":0}`
	response := client.PerformRequest("POST", "/api/collections/", strings.NewReader(jsonBody), map[string]string{
		"Content-Type": "application/json",
	})

	assert.Equal(t, http.StatusBadRequest, response.Code)

	var body map[string]string
	json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, "Invalid input", body["error"])
}

func TestDeleteQuantityFromCollection_Success(t *testing.T) {
	client := api_clients.NewTestClient(true)
	userID := client.User.ID

	card := factories.CardFactory("Blue-Eyes White Dragon")

	factories.UserCardFactory(userID, card.ID, 3)

	payload := map[string]interface{}{
		"quantity": 2,
	}

	url := fmt.Sprintf("/api/collections/%d", card.ID)

	response := client.PerformRequest("DELETE", url, payload, nil)

	assert.Equal(t, http.StatusOK, response.Code)

	var body map[string]string
	json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, "Card quantity updated or removed successfully", body["message"])
}

func TestDeleteQuantityFromCollection_InvalidCardID(t *testing.T) {
	client := api_clients.NewTestClient(true)

	jsonBody := `{"quantity":2}`
	response := client.PerformRequest("DELETE", "/api/collections/abc", strings.NewReader(jsonBody), map[string]string{
		"Content-Type": "application/json",
	})

	assert.Equal(t, http.StatusBadRequest, response.Code)

	var body map[string]string
	json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, "Invalid card ID", body["error"])
}

func TestDeleteQuantityFromCollection_InvalidQuantity(t *testing.T) {
	client := api_clients.NewTestClient(true)

	jsonBody := `{"quantity":0}`
	response := client.PerformRequest("DELETE", "/api/collections/1", strings.NewReader(jsonBody), map[string]string{
		"Content-Type": "application/json",
	})

	assert.Equal(t, http.StatusBadRequest, response.Code)

	var body map[string]string
	json.Unmarshal(response.Body.Bytes(), &body)
	assert.Equal(t, "Invalid quantity", body["error"])
}
