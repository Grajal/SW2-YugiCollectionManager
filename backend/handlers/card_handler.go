package handlers

import (
	"net/http"
	"strconv"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/services"
	"github.com/gin-gonic/gin"
)

func GetOrFetchCard(c *gin.Context) {
	param := c.Param("param")
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

	c.JSON(http.StatusOK, gin.H{"totalCards": total, "cards": cards})
}
