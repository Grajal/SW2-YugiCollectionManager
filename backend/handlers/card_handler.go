package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/gin-gonic/gin"
)

func GetOrFetchCard(c *gin.Context) {
	param := strings.TrimSpace(c.Param("param"))
	if param == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameter (ID or name)"})
		return
	}

	var id int
	var name string

	if parsedID, err := strconv.Atoi(param); err == nil {
		id = parsedID
	} else {
		name = param
	}

	card, err := services.GetOrFetchCardByIDOrName(id, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, card)
}

// GetCards is a Gin handler that retrieves a paginated list of cards from the database.
// If the number of cards in the database is below a defined minimum, it fetches new random cards from the external API.
func GetCards(c *gin.Context) {
	const MinCardCount = 20
	if err := services.EnsureMinimumCards(MinCardCount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ensure minimum card stock"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	cards, err := services.GetCards(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cards"})
		return
	}

	total, err := services.CountAllCards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count all cards"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"totalCards": total, "Some cards": cards})
}

// SearchCards handles GET requests to search cards in the database.
// It filters by name, type, and archetype. If no results are found locally,
// it tries to fetch similar cards from the external YGOProDeck API.
func SearchCards(c *gin.Context) {
	name := c.Query("name")
	cardType := c.Query("type")
	archetype := c.Query("archetype")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	cards, err := services.GetFilteredCards(name, cardType, archetype, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cards from database"})
		return
	}

	if len(cards) == 0 && name != "" {
		fetchedCards, err := services.FetchAndStoreCardsByName(name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cards from external API"})
			return
		}
		cards = fetchedCards
	}

	total, err := services.CountFilteredCards(name, cardType, archetype)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count cards"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"totalCards": total,
		"cards":      cards,
	})
}
