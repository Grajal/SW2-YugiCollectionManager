package services

import (
	"errors"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"gorm.io/gorm"
)

func GetCollectionByUserID(userID uint) ([]models.UserCard, error) {
	var userCards []models.UserCard

	if err := database.DB.Preload("Card").Where("user_id = ?", userID).Find(&userCards).Error; err != nil {
		return nil, err
	}

	return userCards, nil
}

func AddCardToCollection(userID uint, cardID uint, quantity int) error {
	var userCard models.UserCard

	err := database.DB.Where("user_id = ? AND card_id = ?", userID, cardID).First(&userCard).Error
	if err == nil {
		// Card already exists in the collection, update the quantity
		userCard.Quantity += quantity
		return database.DB.Save(&userCard).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newUserCard := models.UserCard{
			UserID:   userID,
			CardID:   cardID,
			Quantity: quantity,
		}

		return database.DB.Create(&newUserCard).Error
	}

	return err
}
