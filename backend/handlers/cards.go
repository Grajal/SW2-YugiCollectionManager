package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/gin-gonic/gin"
)

// Card represents a single card object from the API response
type Card struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	HumanReadableType string `json:"humanReadableCardType"`
	FrameType         string `json:"frameType"`
	Desc              string `json:"desc"`
	Race              string `json:"race"`
	Atk               int    `json:"atk"`
	Def               int    `json:"def"`
	Level             int    `json:"level"`
	Attribute         string `json:"attribute"`
	Archetype         string `json:"archetype"`
}

// ApiResponse represents the root structure of the YGOProDeck API response
type ApiResponse struct {
	Data []Card `json:"data"`
}

// GetCardByName retrieves a single card by name from the YGOProDeck API
func GetCardByName(cardName string) (*Card, error) {
	baseURL := "https://db.ygoprodeck.com/api/v7/cardinfo.php"
	params := url.Values{}
	params.Add("name", cardName)

	resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp ApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(apiResp.Data) == 0 {
		return nil, fmt.Errorf("no card found with name '%s'", cardName)
	}

	return &apiResp.Data[0], nil
}

func GetNewCard(c *gin.Context) {
	cardName := c.Query("name")
	if cardName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'name' query parameter"})
		return
	}

	cardData, err := GetCardByName(cardName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if the card already exists
	exists, err := database.CheckIfCardExists(uint(cardData.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database check failed", "details": err.Error()})
		return
	}

	if exists {
		// Card already exists
		c.JSON(http.StatusConflict, gin.H{"error": "Card already exists"})
		return
	}

	// Check if it's monster or spell card

	// Convert API card to your DB model
	card := models.Card{
		ID:        uint(cardData.ID),
		Name:      cardData.Name,
		Type:      cardData.Type,
		FrameType: cardData.FrameType,
		Desc:      cardData.Desc,
	}

	// Save to database
	if err := database.DB.Create(&card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert card into database", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, card)
}
