package api_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	api_clients "github.com/Grajal/SW2-YugiCollectionManager/backend/tests/clients"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/tests/factories"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetCardNotInAPI(t *testing.T) {
	client := api_clients.NewTestClient(true)

	// Use a card name that you know does NOT exist in external API
	response := client.PerformRequest("GET", "/api/cards/NonExistentCardName", nil, nil)

	// Expect 404 because card not found anywhere
	assert.Equal(t, http.StatusNotFound, response.Code)

	// Assert error message returned
	var body map[string]string
	err := json.Unmarshal(response.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Contains(t, body["error"], "card not found")
}

func TestGetCardMissingParameter(t *testing.T) {
	client := api_clients.NewTestClient(true)

	// Here assuming /api/cards/ expects a card name or id parameter, so calling root /api/cards/
	response := client.PerformRequest("GET", "/api/cards/ ", nil, nil)

	// Expect bad request (400)
	assert.Equal(t, http.StatusBadRequest, response.Code)

	// Assert error message returned
	var body map[string]string
	err := json.Unmarshal(response.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Contains(t, body["error"], "missing parameter")
}

func TestGetCardAlreadyInDb(t *testing.T) {
	client := api_clients.NewTestClient(true)

	// Setup: Insert a card directly into DB (assuming you have a helper or factory)
	cardInDb := models.Card{
		Name: "Blue-Eyes White Dragon",
		// fill in other fields as needed
	}

	factories.CardFactory("Blue-Eyes White Dragon")

	// Request the card by name
	response := client.PerformRequest("GET", "/api/cards/Blue-Eyes White Dragon", nil, nil)

	// Expect 200 OK and card returned
	assert.Equal(t, http.StatusOK, response.Code)

	var returnedCard models.Card
	err := json.Unmarshal(response.Body.Bytes(), &returnedCard)
	assert.NoError(t, err)
	assert.Equal(t, cardInDb.Name, returnedCard.Name)
}

func TestGetCardNotInDb(t *testing.T) {
	// Initialize the test client with authentication (creates a user and token)
	client := api_clients.NewTestClient(true)

	// Override UploadImagetoS3 function (mock for testing)
	utils.UploadImage = func(cardID int, imageURL string) (string, error) {
		return "https://mock-s3/test.jpg", nil
	}

	// Perform a GET request to /items
	response := client.PerformRequest("GET", "/api/cards/Dark Magician", nil, nil)

	// Assert that the response code is 200 OK
	assert.Equal(t, http.StatusOK, response.Code)
	fmt.Println(response.Body.String())

	// Define the expected items response (make sure the structure matches the JSON output)
	expectedBody := models.Card{}

	// Unmarshal the response to a slice of items and assert it matches the expected data
	var actualItems models.Card
	err := json.Unmarshal(response.Body.Bytes(), &actualItems)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
		t.Logf("Received non-200 response: %d, body: %s", response.Code, response.Body.String())
	}

	// Assert that the returned items match what was created
	assert.Equal(t, expectedBody, actualItems)
}
