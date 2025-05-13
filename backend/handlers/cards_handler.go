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
