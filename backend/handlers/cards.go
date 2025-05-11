package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/gin-gonic/gin"
)

const apiBaseURL = "https://db.ygoprodeck.com/api/v7/cardinfo.php"

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

// HttpClient for external API requests, making it easier to test and reuse
var httpClient = &http.Client{}

// GetCardByName retrieves a single card by name from the YGOProDeck API
func GetCardByName(cardName string) (*Card, error) {
	// Prepare the query parameters
	params := url.Values{}
	params.Add("name", cardName)

	// Make the API request
	resp, err := httpClient.Get(fmt.Sprintf("%s?%s", apiBaseURL, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response
	var apiResp ApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Return error if no card is found
	if len(apiResp.Data) == 0 {
		return nil, fmt.Errorf("no card found with name '%s'", cardName)
	}

	// Return the first card found
	return &apiResp.Data[0], nil
}

// CreateCard handles creating a new card in the database and determines the card's type
func CreateCard(card *Card) (*models.Card, error) {

	// Check if card already exists in the database
	exists, err := services.CheckIfCardExists(uint(card.ID))

	if err != nil {
		return nil, fmt.Errorf("database check failed: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("card already exists")
	}

	// Convert API card to DB model
	dbCard := models.Card{
		ID:        uint(card.ID),
		Name:      card.Name,
		Type:      card.Type, // Store the type of the card
		FrameType: card.FrameType,
		Desc:      card.Desc,
	}

	// Begin a transaction to ensure card and its specific type model are created together
	tx := database.DB.Begin()

	// Save the main card record
	if err := tx.Create(&dbCard).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to insert card into database: %w", err)
	}

	// Depending on the card type, create the associated model (MonsterCard or SpellTrapCard)
	if strings.Contains(card.Type, "Monster") {

		// Create a MonsterCard if the card is a Monster
		monsterCard := models.MonsterCard{
			CardID:    dbCard.ID,
			Atk:       card.Atk,
			Def:       card.Def,
			Level:     card.Level,
			Attribute: card.Attribute,
			Race:      card.Race,
		}
		if err := tx.Create(&monsterCard).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to insert monster card into database: %w", err)
		}

	} else if card.Type == "Spell Card" || card.Type == "Trap Card" {

		// Create a SpellTrapCard if the card is a Spell or Trap
		spellTrapCard := models.SpellTrapCard{
			CardID: dbCard.ID,
			Type:   card.Type,
		}
		if err := tx.Create(&spellTrapCard).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to insert spell/trap card into database: %w", err)
		}

	} else {
		// Handle unknown card types if necessary
		tx.Rollback()
		return nil, fmt.Errorf("unknown card type: %s", card.Type)
	}

	// Commit the transaction if everything is successful
	tx.Commit()

	return &dbCard, nil
}

// GetNewCard handles the creation of a new card from the YGOProDeck API
func GetNewCard(c *gin.Context) {
	cardName := c.Query("name")
	if cardName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'name' query parameter"})
		return
	}

	// Get card data from the API
	cardData, err := GetCardByName(cardName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create the card in the database
	createdCard, err := CreateCard(cardData)
	if err != nil {
		// Handle error from database or if card already exists
		if err.Error() == "card already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Respond with the created card
	c.JSON(http.StatusOK, createdCard)
}
