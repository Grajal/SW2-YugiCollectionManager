package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/repository"
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

// DeckCardService defines operations related to managing cards within a deck.
type DeckCardService interface {
	AddCardToDeck(userID, deckID uint, card *models.Card, quantity int) error
	RemoveCardFromDeck(userID, deckID, cardID uint, quantity int) error
}

type deckCardService struct {
	repo repository.DeckCardRepository
}

// NewDeckCardService creates a new instance of deckCardService.
func NewDeckCardService(repo repository.DeckCardRepository) DeckCardService {
	return &deckCardService{
		repo: repo,
	}
}

// AddCardToDeck adds a given quantity of a card to the deck,
// validating deck size and card copy constraints.
func (s *deckCardService) AddCardToDeck(userID, deckID uint, card *models.Card, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("invalid quantity: must be greater than 0")
	}

	zone := GetZoneFromCard(card)

	// Validaciones
	if err := ValidateDeckCardCount(deckID, zone, quantity); err != nil {
		return fmt.Errorf("deck size validation failed: %w", err)
	}
	if err := ValidateCardCopyLimit(deckID, card.ID, quantity); err != nil {
		return fmt.Errorf("copy limit validation failed: %w", err)
	}

	// Persistencia
	if err := s.repo.AddCardToDeck(deckID, card.ID, quantity, zone); err != nil {
		return fmt.Errorf("failed to add card to deck: %w", err)
	}

	return nil
}

// RemoveCardFromDeck removes a quantity of a card from a deck.
// Deletes the card if the resulting quantity is zero or less.
func (s *deckCardService) RemoveCardFromDeck(userID, deckID, cardID uint, quantity int) error {
	entry, err := s.repo.GetDeckCard(deckID, cardID)
	if err != nil {
		return fmt.Errorf("card not found in deck: %w", err)
	}

	if quantity >= entry.Quantity {
		if err := s.repo.DeleteDeckCard(entry); err != nil {
			return fmt.Errorf("failed to delete card: %w", err)
		}
	} else {
		entry.Quantity -= quantity
		if err := s.repo.UpdateDeckCardQuantity(entry); err != nil {
			return fmt.Errorf("failed to update quantity: %w", err)
		}
	}

	return nil
}

// GetZoneFromCard determines whether a card belongs to the "main" or "extra" zone
// based on its frame type.
func GetZoneFromCard(card *models.Card) string {
	switch card.FrameType {
	case "link", "xyz", "fusion", "synchro":
		return "extra"
	default:
		return "main"
	}
}

// ValidateDeckCardCount checks if adding a number of cards to the given zone
// would exceed the allowed deck size.
func ValidateDeckCardCount(deckID uint, zone string, newCards int) error {
	var total sql.NullInt64
	query := database.DB.Model(&models.DeckCard{}).Where("deck_id = ? AND zone = ?", deckID, zone).Select("SUM(quantity)").Scan(&total)
	if query.Error != nil {
		return query.Error
	}

	sum := int64(0)
	if total.Valid {
		sum = total.Int64
	}

	if zone == "extra" && sum+int64(newCards) > 20 {
		return ErrExtraDeckLimitReached
	}
	if zone == "main" && sum+int64(newCards) > 60 {
		return ErrDeckLimitReached
	}
	return nil
}

// ValidateCardCopyLimit ensures that adding new copies of a card
// does not exceed the maximum allowed per deck.
func ValidateCardCopyLimit(deckID, cardID uint, quantityToAdd int) error {
	var existing models.DeckCard
	err := database.DB.Where("deck_id = ? AND card_id = ?", deckID, cardID).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	existingQty := 0
	if existing.DeckID != 0 && existing.CardID != 0 {
		existingQty = existing.Quantity
	}

	if existingQty+quantityToAdd > 3 {
		return ErrCardCopyLimitExceeded
	}

	return nil
}
