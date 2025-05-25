package api_tests

import (
	"encoding/json"
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
	assert.Contains(t, body["error"], "Missing parameter")
}

func TestGetCardAlreadyInDb(t *testing.T) {
	client := api_clients.NewTestClient(true)

	// Setup: Insert a card using the factory
	expectedCard := factories.CardFactory("Blue-Eyes White Dragon")

	// Perform API request to fetch the card by name
	response := client.PerformRequest("GET", "/api/cards/Blue-Eyes White Dragon", nil, nil)

	// Assert response status
	assert.Equal(t, http.StatusOK, response.Code)

	// Unmarshal response body
	var returnedCard models.Card
	err := json.Unmarshal(response.Body.Bytes(), &returnedCard)
	assert.NoError(t, err)

	// Assert key fields match
	assert.Equal(t, expectedCard.ID, returnedCard.ID)
	assert.Equal(t, expectedCard.CardYGOID, returnedCard.CardYGOID)
	assert.Equal(t, expectedCard.Name, returnedCard.Name)
	assert.Equal(t, expectedCard.Desc, returnedCard.Desc)
	assert.Equal(t, expectedCard.Type, returnedCard.Type)
	assert.Equal(t, expectedCard.ImageURL, returnedCard.ImageURL)
	assert.Equal(t, expectedCard.FrameType, returnedCard.FrameType)

	// Optional: if you later extend CardFactory to include MonsterCard, you can assert that too
	// assert.NotNil(t, returnedCard.MonsterCard)
	// assert.Equal(t, expectedCard.MonsterCard.Atk, returnedCard.MonsterCard.Atk)
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

	// Unmarshal the response to a slice of items and assert it matches the expected data
	var actualCard models.Card
	err := json.Unmarshal(response.Body.Bytes(), &actualCard)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
		t.Logf("Received non-200 response: %d, body: %s", response.Code, response.Body.String())
	}

	// Assert that the returned items match what was created
	assert.Equal(t, 46986414, actualCard.CardYGOID)
	assert.Equal(t, "Dark Magician", actualCard.Name)
	assert.Equal(t, "Normal Monster", actualCard.Type)
	assert.Equal(t, "normal", actualCard.FrameType)
	assert.Equal(t, "''The ultimate wizard in terms of attack and defense.''", actualCard.Desc)

	assert.NotNil(t, actualCard.MonsterCard)
	assert.Equal(t, 2500, actualCard.MonsterCard.Atk)
	assert.Equal(t, 2100, actualCard.MonsterCard.Def)
	assert.Equal(t, 7, actualCard.MonsterCard.Level)
	assert.Equal(t, "DARK", actualCard.MonsterCard.Attribute)
	assert.Equal(t, "Spellcaster", actualCard.MonsterCard.Race)

}
