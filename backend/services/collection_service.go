package services

import (
	"errors"
	"fmt"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"gorm.io/gorm"
)

// GetCollectionByUserID retrieves the full collection of cards for a user,
// including all subtypes of the cards (Monster, Spell/Trap, Link, Pendulum).
func GetCollectionByUserID(userID uint) ([]models.UserCard, error) {
	var userCards []models.UserCard

	err := database.DB.
		Preload("Card").
		Preload("Card.MonsterCard").
		Preload("Card.SpellTrapCard").
		Preload("Card.LinkMonsterCard").
		Preload("Card.PendulumMonsterCard").
		Where("user_id = ?", userID).
		Find(&userCards).Error

	if err != nil {
		return nil, err
	}

	return userCards, nil
}

// AddCardToCollection adds a card to the user's collection or updates the quantity if the card already exists.
// It first checks if a UserCard record exists for the given user ID and card ID.
// If the record exists, it increments the quantity and saves the updated record.
// If the record does not exist, it creates a new UserCard record with the specified quantity.
// Returns an error if the database operation fails.
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

// DeleteCardFromCollection removes a card from the user's collection.
// It deletes the UserCard record that matches the given user ID and card ID.
// Returns an error if the database operation fails.
func DeleteCardFromCollection(userID, cardID uint) error {
	return database.DB.Where("user_id = ? AND card_id = ?", userID, cardID).Delete(&models.UserCard{}).Error
}

func DeleteQuantityCardFromCollection(userID, cardID uint, quantityToRemove int) error {
	var userCard models.UserCard

	err := database.DB.Where("user_id= ? AND card_id = ?", userID, cardID).First(&userCard).Error
	if err != nil {
		return fmt.Errorf("card not found in collection: %w", err)
	}

	if userCard.Quantity > quantityToRemove {
		userCard.Quantity -= quantityToRemove
		return database.DB.Save(&userCard).Error
	}

	return database.DB.Delete(&userCard).Error
}
