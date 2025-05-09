package services

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
)

func GetCollectionByUserID(userID uint) ([]models.UserCard, error) {
	var userCards []models.UserCard

	if err := database.DB.Preload("Card").Where("user_id = ?", userID).Find(&userCards).Error; err != nil {
		return nil, err
	}

	return userCards, nil
}
