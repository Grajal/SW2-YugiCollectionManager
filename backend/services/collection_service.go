package services

import (
	"errors"
	"fmt"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/repository"
	"gorm.io/gorm"
)

type CollectionService interface {
	GetUserCollection(userID uint) ([]models.UserCard, error)
	AddCardToCollection(userID uint, cardID uint, quantity int) error
	RemoveCardFromCollection(userID, cardID uint) error
	DecreaseCardQuantity(userID, cardID uint, quantityToRemove int) error
}

type collectionService struct {
	repo repository.CollectionRepository
}

func NewCollectionService(repo repository.CollectionRepository) CollectionService {
	return &collectionService{repo: repo}
}

func (s *collectionService) GetUserCollection(userID uint) ([]models.UserCard, error) {
	collection, err := s.repo.GetUserCollection(userID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch collection for user %d: %w", userID, err)
	}
	return collection, nil
}

func (s *collectionService) AddCardToCollection(userID uint, cardID uint, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	err := s.repo.AddCardToCollection(userID, cardID, quantity)
	if err != nil {
		return fmt.Errorf("failed to add card %d to user %d's collection: %w", cardID, userID, err)
	}
	return nil
}

func (s *collectionService) RemoveCardFromCollection(userID, cardID uint) error {
	err := s.repo.RemoveCardFromCollection(userID, cardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("card %d not found in user %d's collection", cardID, userID)
		}
		return fmt.Errorf("failed to remove card %d from user %d's collection: %w", cardID, userID, err)
	}
	return nil
}

func (s *collectionService) DecreaseCardQuantity(userID, cardID uint, quantityToRemove int) error {
	if quantityToRemove <= 0 {
		return fmt.Errorf("quantity to remove must be greater than zero")
	}

	err := s.repo.DecreaseCardQuantity(userID, cardID, quantityToRemove)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("card %d not found in user %d's collection", cardID, userID)
		}
		return fmt.Errorf("failed to decrease quantity of card %d in user %d's collection: %w", cardID, userID, err)
	}
	return nil
}
