package repository

import (
	"errors"
	"fmt"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"gorm.io/gorm"
)

type DeckRepository interface {
	CountByUserID(userID uint) (int64, error)
	ExistsByName(userID uint, name string) (bool, error)
	Create(deck *models.Deck) error
	FindByUserID(userID uint) ([]models.Deck, error)
	FindByIDAndUserID(deckID, userID uint) (*models.Deck, error)
	DeleteByIDAndUserID(deckID, userID uint) error
	FindDeckCards(deckID, userID uint) ([]models.DeckCard, error)
}

type deckRepository struct {
	db *gorm.DB
}

func NewDeckRepository() DeckRepository {
	return &deckRepository{
		db: database.DB,
	}
}

// Count decks created by a given user
func (r *deckRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Deck{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// Check if a deck with the same name exists for the user
func (r *deckRepository) ExistsByName(userID uint, name string) (bool, error) {
	var deck models.Deck
	err := r.db.Where("user_id = ? AND name = ?", userID, name).First(&deck).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

// Create a new deck
func (r *deckRepository) Create(deck *models.Deck) error {
	return r.db.Create(deck).Error
}

// Find all decks belonging to a user, preloading deck cards and their types
func (r *deckRepository) FindByUserID(userID uint) ([]models.Deck, error) {
	var decks []models.Deck
	err := r.db.Where("user_id = ?", userID).
		Preload("DeckCards").
		Preload("DeckCards.Card").
		Preload("DeckCards.Card.MonsterCard").
		Preload("DeckCards.Card.SpellTrapCard").
		Preload("DeckCards.Card.LinkMonsterCard").
		Preload("DeckCards.Card.PendulumMonsterCard").
		Find(&decks).Error
	return decks, err
}

// Find a deck by ID and user ID
func (r *deckRepository) FindByIDAndUserID(deckID, userID uint) (*models.Deck, error) {
	var deck models.Deck
	err := r.db.First(&deck, "id = ? AND user_id = ?", deckID, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("deck not found")
	}
	return &deck, err
}

// Delete a deck by ID and user ID
func (r *deckRepository) DeleteByIDAndUserID(deckID, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", deckID, userID).Delete(&models.Deck{}).Error
}

// Get all deck cards for a specific deck
func (r *deckRepository) FindDeckCards(deckID, userID uint) ([]models.DeckCard, error) {
	var deck models.Deck
	err := r.db.Where("id = ? AND user_id = ?", deckID, userID).First(&deck).Error
	if err != nil {
		return nil, fmt.Errorf("deck not found or access denied")
	}

	var deckCards []models.DeckCard
	err = r.db.Where("deck_id = ?", deckID).
		Preload("Card").
		Preload("Card.MonsterCard").
		Preload("Card.SpellTrapCard").
		Preload("Card.LinkMonsterCard").
		Preload("Card.PendulumMonsterCard").
		Find(&deckCards).Error

	return deckCards, err
}
