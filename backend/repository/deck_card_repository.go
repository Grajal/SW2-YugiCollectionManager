package repository

import (
	"errors"
	"fmt"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"gorm.io/gorm"
)

// DeckCardRepository defines the interface for operations on deck-card relations.
type DeckCardRepository interface {
	AddCardToDeck(deckID, cardID uint, quantity int, zone string) error
	GetDeckCard(deckID, cardID uint) (*models.DeckCard, error)
	UpdateDeckCardQuantity(card *models.DeckCard) error
	DeleteDeckCard(card *models.DeckCard) error
}

type deckCardRepository struct {
	db *gorm.DB
}

// NewDeckCardRepository creates a new instance of deckCardRepository using the default DB.
func NewDeckCardRepository() DeckCardRepository {
	return &deckCardRepository{
		db: database.DB,
	}
}

// AddCardToDeck adds a card to a deck, or updates the quantity if the card already exists.
func (r *deckCardRepository) AddCardToDeck(deckID, cardID uint, quantity int, zone string) error {
	existing, err := r.GetDeckCard(deckID, cardID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error querying existing deck card: %w", err)
	}

	if existing != nil {
		existing.Quantity += quantity
		return r.UpdateDeckCardQuantity(existing)
	}

	newEntry := &models.DeckCard{
		DeckID:   deckID,
		CardID:   cardID,
		Quantity: quantity,
		Zone:     zone,
	}
	return r.db.Create(newEntry).Error
}

// GetDeckCard retrieves the DeckCard entry (with full card data) for a given deck and card ID.
func (r *deckCardRepository) GetDeckCard(deckID, cardID uint) (*models.DeckCard, error) {
	var card models.DeckCard
	err := r.db.Where("deck_id = ? AND card_id = ?", deckID, cardID).
		Preload("Card").
		Preload("Card.MonsterCard").
		Preload("Card.SpellTrapCard").
		Preload("Card.LinkMonsterCard").
		Preload("Card.PendulumMonsterCard").
		First(&card).Error

	if err != nil {
		return nil, err
	}
	return &card, nil
}

// UpdateDeckCardQuantity updates the quantity of an existing DeckCard entry.
func (r *deckCardRepository) UpdateDeckCardQuantity(card *models.DeckCard) error {
	return r.db.Save(card).Error
}

// DeleteDeckCard removes a DeckCard entry from the database.
func (r *deckCardRepository) DeleteDeckCard(card *models.DeckCard) error {
	return r.db.Delete(card).Error
}
