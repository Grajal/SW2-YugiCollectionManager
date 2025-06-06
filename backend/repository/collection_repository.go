package repository

import (
	"errors"
	"fmt"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"gorm.io/gorm"
)

type CollectionRepository interface {
	GetUserCollection(userID uint) ([]models.UserCard, error)
	GetUserCard(userID, cardID uint) (*models.UserCard, error)
	AddCardToCollection(userID uint, cardID uint, quantity int) error
	DecreaseCardQuantity(userID, cardID uint, quantityToRemove int) error
}

type collectionRepository struct {
	db *gorm.DB
}

func NewCollectionRepository() CollectionRepository {
	return &collectionRepository{
		db: database.DB,
	}
}

func NewCollectionRepositoryWithDB(db *gorm.DB) CollectionRepository {
	return &collectionRepository{
		db: db,
	}
}

func (r *collectionRepository) GetUserCollection(userID uint) ([]models.UserCard, error) {
	var userCards []models.UserCard

	err := r.db.
		Preload("Card").
		Preload("Card.MonsterCard").
		Preload("Card.SpellTrapCard").
		Preload("Card.LinkMonsterCard").
		Preload("Card.PendulumMonsterCard").
		Where("user_id = ?", userID).
		Find(&userCards).Error

	return userCards, err
}

func (r *collectionRepository) GetUserCard(userID, cardID uint) (*models.UserCard, error) {
	var userCard models.UserCard
	err := r.db.Preload("Card").
		Preload("Card.MonsterCard").
		Preload("Card.SpellTrapCard").
		Preload("Card.LinkMonsterCard").
		Preload("Card.PendulumMonsterCard").
		Preload("User").
		First(&userCard, "user_id = ? AND card_id = ?", userID, cardID).Error

	return &userCard, err
}

func (r *collectionRepository) AddCardToCollection(userID uint, cardID uint, quantity int) error {
	var userCard models.UserCard

	err := r.db.Where("user_id = ? AND card_id = ?", userID, cardID).First(&userCard).Error
	if err == nil {
		userCard.Quantity += quantity
		return r.db.Save(&userCard).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newUserCard := models.UserCard{
			UserID:   userID,
			CardID:   cardID,
			Quantity: quantity,
		}
		return r.db.Create(&newUserCard).Error
	}

	return err
}

func (r *collectionRepository) RemoveCard(userID, cardID uint) error {
	var userCard models.UserCard
	if err := r.db.Where("user_id = ? AND card_id = ?", userID, cardID).First(&userCard).Error; err != nil {
		return fmt.Errorf("user card not found: %w", err)
	}
	return r.db.Delete(&userCard).Error
}

func (r *collectionRepository) DecreaseCardQuantity(userID, cardID uint, quantityToRemove int) error {
	var userCard models.UserCard

	err := r.db.Where("user_id = ? AND card_id = ?", userID, cardID).First(&userCard).Error
	if err != nil {
		return fmt.Errorf("card not found in collection: %w", err)
	}

	if userCard.Quantity > quantityToRemove {
		userCard.Quantity -= quantityToRemove
		return r.db.Save(&userCard).Error
	}

	return r.db.Delete(&userCard).Error
}
