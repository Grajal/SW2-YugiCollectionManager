package services

import (
	"errors"
	"fmt"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
)

var ErrDeckAlreadyExists = errors.New("deck with the same name already exists")
var ErrMaximumNumberOfDecks = errors.New("maximum number of decks reached")

func CreateDeck(userID uint, name, description string) (*models.Deck, error) {
	var count int64
	if err := database.DB.Model(&models.Deck{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to check existing decks: %w", err)
	}

	if count >= 10 {
		return nil, ErrMaximumNumberOfDecks
	}

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
