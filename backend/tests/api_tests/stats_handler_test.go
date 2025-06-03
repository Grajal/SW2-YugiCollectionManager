package api_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	api_clients "github.com/Grajal/SW2-YugiCollectionManager/backend/tests/clients"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/tests/factories"

	"github.com/stretchr/testify/assert"
)

func TestGetCollectionStats_Success(t *testing.T) {
	client := api_clients.NewTestClient(true) // Authenticated

	// Optionally set up cards in the user's collection here using factories

	response := client.PerformRequest("GET", "/api/stats/collection", nil, nil)

	assert.Equal(t, http.StatusOK, response.Code)

	var stats map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &stats)
	assert.NoError(t, err)

	// Basic checks - customize depending on actual stat fields
	assert.Contains(t, stats, "total_cards")
}

func TestGetCollectionStats_Unauthorized(t *testing.T) {
	client := api_clients.NewTestClient(false) // Unauthenticated

	response := client.PerformRequest("GET", "/api/stats/collection", nil, nil)

	assert.Equal(t, http.StatusUnauthorized, response.Code)

	var body map[string]string
	err := json.Unmarshal(response.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Contains(t, body["error"], "Authorization token not provided")
}

func TestGetDeckStats_Success(t *testing.T) {
	client := api_clients.NewTestClient(true)

	deck := factories.DeckFactory(&client.User, true) // with cards in "Main" zone

	resp := client.PerformRequest("GET", fmt.Sprintf("/api/stats/deck/%d", deck.ID), nil, nil)

	assert.Equal(t, http.StatusOK, resp.Code)

	var stats map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &stats)
	assert.NoError(t, err)

	assert.Contains(t, stats, "total_cards")
}

func TestGetDeckStats_InvalidDeckID(t *testing.T) {
	client := api_clients.NewTestClient(true)

	response := client.PerformRequest("GET", "/api/stats/deck/invalid", nil, nil)

	assert.Equal(t, http.StatusBadRequest, response.Code)

	var body map[string]string
	err := json.Unmarshal(response.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Contains(t, body["error"], "Invalid deck ID")
}

func TestGetDeckStats_Unauthorized(t *testing.T) {
	client := api_clients.NewTestClient(false)

	response := client.PerformRequest("GET", "/api/stats/deck/1", nil, nil)

	assert.Equal(t, http.StatusUnauthorized, response.Code)

	var body map[string]string
	err := json.Unmarshal(response.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Contains(t, body["error"], "Authorization token not provided")
}

func TestGetDeckStats_DeckNotFound(t *testing.T) {
	client := api_clients.NewTestClient(true)

	// Assuming 99999 is a non-existent deck ID for this user
	response := client.PerformRequest("GET", "/api/stats/deck/99999", nil, nil)

	assert.Equal(t, http.StatusInternalServerError, response.Code)

	var body map[string]string
	err := json.Unmarshal(response.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Contains(t, body["error"], "Could not fetch deck")
}
