package services

import (
	"errors"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
)

var ErrDeckAlreadyExists = errors.New("deck with the same name already exists")

func CreateDeck(userID uint, name, description string) (*models.Deck, error) {
	var existing models.Deck
	if err := database.DB.Where("user_id = ? AND name = ?", userID, name).First(&existing).Error; err == nil {
		return nil, ErrDeckAlreadyExists
	}

	deck := models.Deck{
		UserID:      userID,
		Name:        name,
		Description: description,
	}

	if err := database.DB.Create(&deck).Error; err != nil {
		return nil, err
	}

	return &deck, nil
}
