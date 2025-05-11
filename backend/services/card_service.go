package services

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"gorm.io/gorm"
)

// CheckIfCardExists checks if a card with a given ID exists in the database
func CheckIfCardExists(cardID uint) (bool, error) {
	var card models.Card
	result := database.DB.First(&card, "id = ?", cardID)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
