package services

import (
	"errors"
	"fmt"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"gorm.io/gorm"
)

var ErrDeckAlreadyExists = errors.New("deck with the same name already exists")
var ErrMaximumNumberOfDecks = errors.New("maximum number of decks reached")
var ErrDeckNotFound = errors.New("deck not found")

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

func getDeckByIDAndUserID(deckID, userID uint) (*models.Deck, error) {
	var deck models.Deck
	err := database.DB.First(&deck, "id = ? AND user_id = ?", deckID, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeckNotFound
		}
		return nil, err
	}
	return &deck, nil
}

func GetDecksByUserID(userID uint) ([]models.Deck, error) {
	var decks []models.Deck
	err := database.DB.Where("user_id = ?", userID).Find(&decks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch decks: %w", err)
	}

	return decks, nil
}

func DeleteDeck(deckID uint, userID uint) error {
	var deck models.Deck
	err := database.DB.Where("id = ? AND user_id = ?", deckID, userID).First(&deck).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDeckNotFound
	}
	if err != nil {
		return err
	}

	return database.DB.Delete(&deck).Error
}

func GetCardsByDeck(userID, deckID uint) ([]models.DeckCard, error) {
	var deck models.Deck
	err := database.DB.Where("id = ? AND user_id = ?", deckID, userID).First(&deck).Error
	if err != nil {
		return nil, fmt.Errorf("deck not found or acces denied")
	}

	var deckCards []models.DeckCard
	err = database.DB.Where("deck_id = ?", deckID).Preload("Card").Preload("Card.MonsterCard").Preload("Card.SpellTrapCard").Preload("Card.LinkMonsterCard").Preload("Card.PendulumMonsterCard").Find(&deckCards).Error

	if err != nil {
		return nil, err
	}

	return deckCards, nil
}
