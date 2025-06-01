package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/repository"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/utils"
	"gorm.io/gorm"
)

var ErrDeckAlreadyExists = errors.New("deck with the same name already exists")
var ErrMaximumNumberOfDecks = errors.New("maximum number of decks reached")
var ErrDeckNotFound = errors.New("deck not found")

type DeckService interface {
	CreateDeck(userID uint, name, description string) (*models.Deck, error)
	GetDecksByUserID(userID uint) ([]models.Deck, error)
	DeleteDeck(deckID uint, userID uint) error
	GetCardsByDeck(userID, deckID uint) ([]models.DeckCard, error)
	ExportDeckAsYDK(userID, deckID uint) (string, error)
	ImportDeckFromYDK(userID, deckID uint, file multipart.File) error
	AddCardToDeck(userID, cardID, deckID uint, quantity int) error
	RemoveCardFromDeck(userID, deckID, cardID uint, quantity int) error
}

type deckService struct {
	repo            repository.DeckRepository
	cardService     CardService
	deckCardService DeckCardService
}

func NewDeckService(repo repository.DeckRepository, cardService CardService, deckCardService DeckCardService) DeckService {
	return &deckService{repo, cardService, deckCardService}
}

func (s *deckService) CreateDeck(userID uint, name, description string) (*models.Deck, error) {
	count, err := s.repo.CountByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing decks: %w", err)
	}
	if count >= 10 {
		return nil, ErrMaximumNumberOfDecks
	}

	exists, err := s.repo.ExistsByName(userID, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDeckAlreadyExists
	}

	deck := &models.Deck{
		UserID:      userID,
		Name:        name,
		Description: description,
	}

	if err := s.repo.Create(deck); err != nil {
		return nil, err
	}

	return deck, nil
}

func (s *deckService) GetDecksByUserID(userID uint) ([]models.Deck, error) {
	return s.repo.FindByUserID(userID)
}

func (s *deckService) DeleteDeck(deckID uint, userID uint) error {
	err := s.repo.DeleteByIDAndUserID(deckID, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDeckNotFound
	}
	return err
}

func (s *deckService) GetCardsByDeck(userID, deckID uint) ([]models.DeckCard, error) {
	return s.repo.FindDeckCards(deckID, userID)
}

func (s *deckService) ExportDeckAsYDK(userID, deckID uint) (string, error) {
	deck, err := s.repo.FindByIDAndUserID(deckID, userID)
	if err != nil {
		return "", fmt.Errorf("deck not found or unauthorized: %w", err)
	}

	deckCards, err := s.repo.FindDeckCards(deck.ID, userID)
	if err != nil {
		return "", fmt.Errorf("failed to load cards: %w", err)
	}

	var mainLines, extraLines, sideLines []string
	for _, c := range deckCards {
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

func (s *deckService) ImportDeckFromYDK(userID, deckID uint, file multipart.File) error {
	mainIDs, extraIDs, sideIDs, err := utils.ParseYDK(file)
	if err != nil {
		return fmt.Errorf("error parsing YDK file: %w", err)
	}

	process := func(ids []string) error {
		for _, idStr := range ids {
			cardYGOID, err := strconv.Atoi(idStr)
			if err != nil {
				return fmt.Errorf("invalid card ID %s: %w", idStr, err)
			}

			card, err := s.cardService.GetCardByYGOID(cardYGOID)
			if err != nil {
				return fmt.Errorf("error retrieving card %d: %w", cardYGOID, err)
			}

			if err := s.deckCardService.AddCardToDeck(userID, deckID, card, 1); err != nil {
				return fmt.Errorf("error adding card %d to deck: %w", card.ID, err)
			}
		}
		return nil
	}

	if err := process(mainIDs); err != nil {
		return fmt.Errorf("error processing main cards: %w", err)
	}
	if err := process(extraIDs); err != nil {
		return fmt.Errorf("error processing extra cards: %w", err)
	}
	if err := process(sideIDs); err != nil {
		return fmt.Errorf("error processing side cards: %w", err)
	}

	return nil
}

func (s *deckService) AddCardToDeck(userID, cardID, deckID uint, quantity int) error {
	// Obtener la carta desde el servicio
	card, err := s.cardService.GetCardByID(cardID)
	if err != nil {
		return fmt.Errorf("failed to retrieve card: %w", err)
	}

	// Delegar la lógica de validación y persistencia al servicio de deckCard
	err = s.deckCardService.AddCardToDeck(userID, deckID, card, quantity)
	if err != nil {
		return fmt.Errorf("failed to add card to deck: %w", err)
	}

	return nil
}

func (s *deckService) RemoveCardFromDeck(userID, deckID, cardID uint, quantity int) error {
	return s.deckCardService.RemoveCardFromDeck(userID, deckID, cardID, quantity)
}
