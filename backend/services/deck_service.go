package services

import (
	"errors"
	"fmt"
	"strings"

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
	err := database.DB.Where("user_id = ?", userID).Preload("DeckCards").
		Preload("DeckCards.Card.MonsterCard").
		Preload("DeckCards.Card.SpellTrapCard").
		Preload("DeckCards.Card.LinkMonsterCard").
		Preload("DeckCards.Card.PendulumMonsterCard").
		Preload("DeckCards.Card").Find(&decks).Error
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

func ExportDeckAsYDK(userID, deckID uint) (string, error) {
	deck, err := getDeckByIDAndUserID(deckID, userID)
	if err != nil {
		return "", fmt.Errorf("deck not found or unauthorized: %w", err)
	}

	var cards []models.DeckCard
	err = database.DB.Where("deck_id = ?", deck.ID).Preload("Card").Find(&cards).Error
	if err != nil {
		return "", fmt.Errorf("failed to load cards: %w", err)
	}

	var mainLines, extraLines, sideLines []string
	for _, c := range cards {
		for i := 0; i < c.Quantity; i++ {
			line := fmt.Sprintf("%d", c.Card.CardYGOID)
			switch c.Zone {
			case "main":
				mainLines = append(mainLines, line)
			case "extra":
				extraLines = append(extraLines, line)
			case "side":
				sideLines = append(sideLines, line)
			}
		}
	}

	result := "#main\n" + strings.Join(mainLines, "\n") + "\n#extra\n" + strings.Join(extraLines, "\n") + "\n#side\n" + strings.Join(sideLines, "\n")

	return result, nil
}

// func ImportDeckFromYDK(userID, deckID uint, file multipart.File) error {
// 	mainIDs, extraIDs, sideIDs, err := utils.ParseYDK(file)
// 	if err != nil {
// 		return fmt.Errorf("error parsing YDK file: %w", err)
// 	}

// 	process := func(ids []string) error {
// 		for _, idStr := range ids {
// 			id, err := strconv.Atoi(idStr)
// 			if err != nil {
// 				return fmt.Errorf("invalid card ID %s: %w", idStr, err)
// 			}

// 			card, err := GetOrFetchCardByIDOrName(id, "")
// 			if err != nil {
// 				return fmt.Errorf("error resolving card %d: %w", id, err)
// 			}

// 			_, err = AddCardToDeck(userID, deckID, card.ID, 1)
// 			if err != nil {
// 				return fmt.Errorf("error adding card %d to deck: %w", card.ID, err)
// 			}
// 		}
// 		return nil
// 	}

// 	if err := process(mainIDs); err != nil {
// 		return fmt.Errorf("error processing main cards: %w", err)
// 	}
// 	if err := process(extraIDs); err != nil {
// 		return fmt.Errorf("error processing extra cards: %w", err)
// 	}
// 	if err := process(sideIDs); err != nil {
// 		return fmt.Errorf("error processing side cards: %w", err)
// 	}

// 	return nil
// }
