package services

import (
	"fmt"
	"log"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/client"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/repository"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/utils"
)

// CardService defines the operations available to manage cards.
// Implementations should be responsible for obtaining, filtering and persisting cards.
type CardService interface {
	GetCardByID(id uint) (*models.Card, error)
	GetCardByYGOID(id int) (*models.Card, error)
	GetCardByName(name string) (*models.Card, error)
	GetCards(limit, offset int) ([]*models.Card, error)
	CountAllCards() (int64, error)
	GetFilteredCards(name, cardType, frameType string, limit, offset int) ([]*models.Card, error)
	CountFilteredCards(name, cardType, frameType string) (int64, error)
}

type cardService struct {
	repo    repository.CardRepository
	factory CardFactory
}

func NewCardService(repo repository.CardRepository, factory CardFactory) CardService {
	return &cardService{repo, factory}
}

// GetCardByID retrieves a card from the database using its internal database ID.
// Returns an error if the card is not found or the database query fails.
func (s *cardService) GetCardByID(id uint) (*models.Card, error) {
	return s.repo.GetByID(id)
}

// GetCardByYGOID retrieves a card by its YGOProDeck ID.
// If the card does not exist in the local database, it attempts to fetch it from the external API,
// uploads the image to S3, builds the card, saves it, and returns the resulting model.
func (s *cardService) GetCardByYGOID(id int) (*models.Card, error) {
	// Paso 1: buscar en base de datos
	card, err := s.repo.GetByYGOProID(id)
	if err == nil && card != nil {
		return card, nil
	}

	// Paso 2: fetch desde API externa
	apiCard, err := client.FetchCardByIDOrName(id, "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch card from API: %w", err)
	}

	if apiCard == nil || apiCard.ImageURL == "" {
		return nil, fmt.Errorf("invalid API response: missing card or image")
	}

	imageURL, err := utils.UploadCardImageToS3(apiCard.ID, apiCard.ImageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image to S3: %w", err)
	}

	card = s.factory.BuildCardFromAPI(apiCard, imageURL)

	if err := s.repo.Create(card); err != nil {
		return nil, fmt.Errorf("failed to save card to database: %w", err)
	}

	return card, nil
}

// GetCardByName retrieves a card by its name.
// If not found in the database, it tries to fetch it from the external API,
// uploads the image to S3, builds the card, saves it, and returns it.
func (s *cardService) GetCardByName(name string) (*models.Card, error) {
	card, err := s.repo.GetByName(name)
	if err == nil && card != nil {
		return card, nil
	}

	apiCard, err := client.FetchCardByIDOrName(0, name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch card from API: %w", err)
	}

	if apiCard == nil || apiCard.ImageURL == "" {
		return nil, fmt.Errorf("invalid API response: missing card or image")
	}

	imageURL, err := utils.UploadCardImageToS3(apiCard.ID, apiCard.ImageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image to S3: %w", err)
	}

	card = s.factory.BuildCardFromAPI(apiCard, imageURL)

	if err := s.repo.Create(card); err != nil {
		return nil, fmt.Errorf("failed to save card to database: %w", err)
	}

	return card, nil
}

// GetCards retrieves all cards from the database using pagination.
// If the total number of cards in the database is below a defined threshold,
// it attempts to fetch a set of random cards from the external API and store them locally.
func (s *cardService) GetCards(limit, offset int) ([]*models.Card, error) {
	const minCardThreshold = 50
	const cardsToFetch = 20

	total, err := s.repo.CountAll()
	if err != nil {
		return nil, fmt.Errorf("failed to count cards: %w", err)
	}

	if total < minCardThreshold {
		randomCards, err := client.FetchRandomCards(cardsToFetch)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch random cards: %w", err)
		}

		for _, apiCard := range randomCards {
			if len(apiCard.CardImages) == 0 || apiCard.CardImages[0].ImageURL == "" {
				continue
			}

			imageURL, err := utils.UploadCardImageToS3(apiCard.ID, apiCard.CardImages[0].ImageURL)
			if err != nil {
				log.Printf("error uploading image for card %s: %v", apiCard.Name, err)
				continue
			}

			card := s.factory.BuildCardFromAPI(&apiCard, imageURL)
			if err := s.repo.Create(card); err != nil {
				log.Printf("error saving card %s: %v", card.Name, err)
			}
		}
	}

	cards, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Card, 0, len(cards))
	for i := range cards {
		result = append(result, &cards[i])
	}
	return result, nil
}

// CountAllCards returns the total number of cards stored in the database.
func (s *cardService) CountAllCards() (int64, error) {
	return s.repo.CountAll()
}

// GetFilteredCards returns cards that match the provided name, type, and frameType filters.
// If no cards are found in the local database and a name is provided,
// it attempts to fetch matching cards from the external API and stores them locally.
func (s *cardService) GetFilteredCards(name, cardType, frameType string, limit, offset int) ([]*models.Card, error) {
	cards, err := s.repo.GetFiltered(name, cardType, frameType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to filter cards from database: %w", err)
	}

	if len(cards) > 0 {
		result := make([]*models.Card, 0, len(cards))
		for i := range cards {
			result = append(result, &cards[i])
		}
		return result, nil
	}

	if name != "" {
		apiCards, err := client.FetchCardsByName(name)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch cards from external API: %w", err)
		}

		var fetched []*models.Card
		for _, apiCard := range apiCards {
			if len(apiCard.CardImages) == 0 || apiCard.CardImages[0].ImageURL == "" {
				continue
			}

			imageURL, err := utils.UploadCardImageToS3(apiCard.ID, apiCard.CardImages[0].ImageURL)
			if err != nil {
				log.Printf("error uploading image for card %s: %v", apiCard.Name, err)
				continue
			}

			card := s.factory.BuildCardFromAPI(&apiCard, imageURL)

			if err := s.repo.Create(card); err != nil {
				log.Printf("error saving card %s: %v", card.Name, err)
				continue
			}

			fetched = append(fetched, card)
		}

		return fetched, nil
	}

	return []*models.Card{}, nil
}

// CountFilteredCards returns the number of cards in the database that match the provided filters.
// Returns 0 and an error if the database query fails.
func (s *cardService) CountFilteredCards(name, cardType, frameType string) (int64, error) {
	count, err := s.repo.CountFiltered(name, cardType, frameType)
	if err != nil {
		return 0, fmt.Errorf("failed to count filtered cards: %w", err)
	}
	return count, nil
}
