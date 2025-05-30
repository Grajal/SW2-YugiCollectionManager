package services

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"gorm.io/gorm"
)

var (
	ErrCardNotFound          = errors.New("card not found")
	ErrCardCopyLimitExceeded = errors.New("too many copies of this card in the deck")
	ErrDeckLimitReached      = errors.New("deck size limit reached")
	ErrExtraDeckLimitReached = errors.New("extra deck size limit reached")
)

const (
	MainDeckMaxSize  = 60
	ExtraDeckMaxSize = 20
	MaxCopiesPerCard = 3
)

func AddCardToDeck(userID uint, deckID uint, cardID uint, quantity int) (*models.DeckCard, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be grater than 0")
	}

	deck, err := getDeckByIDAndUserID(deckID, userID)
	if err != nil {
		return nil, fmt.Errorf("deck not found or unauthorized: %w", err)
	}

	card, err := GetCardByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("card not found: %w", err)
	}

	isExtra := IsExtraDeckCard(card.FrameType)

	if err := ValidateDeckCardCount(deck.ID, isExtra, quantity); err != nil {
		return nil, err
	}

	if err := ValidateCardCopyLimit(deck.ID, card.ID, quantity); err != nil {
		return nil, err
	}

	var existing models.DeckCard
	err = database.DB.Where("deck_id = ? AND card_id = ?", deck.ID, card.ID).First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to query card: %w", err)
	}

	if existing.ID != 0 {
		existing.Quantity += quantity
		if err := database.DB.Save(&existing).Error; err != nil {
			return nil, fmt.Errorf("failed to update card quantity: %w", err)
		}
		return &existing, nil
	}

	newEntry := models.DeckCard{
		DeckID:      deckID,
		CardID:      card.ID,
		Quantity:    quantity,
		IsExtraDeck: isExtra,
	}

	if err := database.DB.Create(&newEntry).Error; err != nil {
		return nil, fmt.Errorf("failed to add card to deck: %w", err)
	}

	return &newEntry, nil
}

func IsExtraDeckCard(frameType string) bool {
	switch strings.ToLower(frameType) {
	case "link", "xyz", "fusion", "synchro":
		return true
	default:
		return false
	}
}

func RemoveCardFromDeck(userID, deckID, cardID uint, quantity int) error {
	deck, err := getDeckByIDAndUserID(deckID, userID)
	if err != nil {
		return fmt.Errorf("deck not found or unathorized: %w", err)
	}

	var entry models.DeckCard
	err = database.DB.Where("deck_id = ? AND card_id = ?", deck.ID, cardID).First(&entry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("card not found in deck")
		}
		return fmt.Errorf("error fetching card from deck: %w", err)
	}

	if quantity >= entry.Quantity {
		if err := database.DB.Delete(&entry).Error; err != nil {
			return fmt.Errorf("failed to delete card from deck: %w", err)
		}
	} else {
		entry.Quantity -= quantity
		if err := database.DB.Save(&entry).Error; err != nil {
			return fmt.Errorf("failed to update card quantity: %w", err)
		}
	}

	return nil
}

func ValidateDeckCardCount(deckID uint, isExtra bool, newCards int) error {
	var total sql.NullInt64
	query := database.DB.Model(&models.DeckCard{}).Where("deck_id = ? AND is_extra_deck = ?", deckID, isExtra).Select("SUM(quantity)").Scan(&total)
	if query.Error != nil {
		return query.Error
	}

	sum := int64(0)
	if total.Valid {
		sum = total.Int64
	}

	if isExtra && sum+int64(newCards) > 20 {
		return ErrExtraDeckLimitReached
	}
	if !isExtra && sum+int64(newCards) > 60 {
		return ErrDeckLimitReached
	}
	return nil
}

func ValidateCardCopyLimit(deckID, cardID uint, quantityToAdd int) error {
	var existing models.DeckCard
	err := database.DB.Where("deck_id = ? AND card_id = ?", deckID, cardID).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	existingQty := 0
	if existing.ID != 0 {
		existingQty = existing.Quantity
	}

	if existingQty+quantityToAdd > 3 {
		return ErrCardCopyLimitExceeded
	}

	return nil
}
