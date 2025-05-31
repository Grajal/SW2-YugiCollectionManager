package repository

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"gorm.io/gorm"
)

// CardRepository defines the interface for accessing and managing Card data in the database.
type CardRepository interface {
	GetByID(id uint) (*models.Card, error)
	GetByYGOProID(id int) (*models.Card, error)
	GetByName(name string) (*models.Card, error)
	GetAll(limit, offset int) ([]models.Card, error)
	CountAll() (int64, error)
	GetFiltered(name, cardType, frameType string, limit, offset int) ([]models.Card, error)
	CountFiltered(name, cardType, frameType string) (int64, error)
	Create(card *models.Card) error
	ExistsByYGOProID(id int) (bool, error)
}

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository() CardRepository {
	return &cardRepository{
		db: database.DB,
	}
}

// GetByID retrieves a Card by its internal database ID, including all related subtypes.
func (r *cardRepository) GetByID(id uint) (*models.Card, error) {
	var card models.Card
	err := r.db.Preload("MonsterCard").
		Preload("SpellTrapCard").
		Preload("LinkMonsterCard").
		Preload("PendulumMonsterCard").
		First(&card, "id = ?", id).Error

	return &card, err
}

// GetByYGOProID retrieves a Card by its YGOProDeck ID, including all related subtypes.
func (r *cardRepository) GetByYGOProID(id int) (*models.Card, error) {
	var card models.Card
	err := r.db.Preload("MonsterCard").
		Preload("SpellTrapCard").
		Preload("LinkMonsterCard").
		Preload("PendulumMonsterCard").
		First(&card, "card_ygo_id = ?", id).Error
	return &card, err
}

// GetByName retrieves a Card by its exact name, case-insensitive
func (r *cardRepository) GetByName(name string) (*models.Card, error) {
	var card models.Card
	err := r.db.Preload("MonsterCard").
		Preload("SpellTrapCard").
		Preload("LinkMonsterCard").
		Preload("PendulumMonsterCard").
		First(&card, "name ILIKE ?", name).Error
	return &card, err
}

// GetAll retrieves a paginated list of Cards, including their subtypes.
func (r *cardRepository) GetAll(limit, offset int) ([]models.Card, error) {
	var cards []models.Card
	err := r.db.Preload("MonsterCard").
		Preload("SpellTrapCard").
		Preload("LinkMonsterCard").
		Preload("PendulumMonsterCard").
		Limit(limit).Offset(offset).
		Find(&cards).Error
	return cards, err
}

// CountAll returns the total number of Cards in the database.
func (r *cardRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.Card{}).Count(&count).Error
	return count, err
}

// GetFiltered retrieves a paginated list of Cards matching the given filters (name, type, frameType).
func (r *cardRepository) GetFiltered(name, cardType, frameType string, limit, offset int) ([]models.Card, error) {
	query := r.db.Model(&models.Card{}).
		Preload("MonsterCard").
		Preload("SpellTrapCard").
		Preload("LinkMonsterCard").
		Preload("PendulumMonsterCard")

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if cardType != "" {
		query = query.Where("type = ?", cardType)
	}
	if frameType != "" {
		query = query.Where("frame_type = ?", frameType)
	}

	var cards []models.Card
	err := query.Limit(limit).Offset(offset).Find(&cards).Error
	return cards, err
}

// CountFiltered returns the number of Cards that match the given filters.
func (r *cardRepository) CountFiltered(name, cardType, frameType string) (int64, error) {
	query := r.db.Model(&models.Card{})
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if cardType != "" {
		query = query.Where("type = ?", cardType)
	}
	if frameType != "" {
		query = query.Where("frame_type = ?", frameType)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// Create saves a new Card and its associated subtype data into the database.
func (r *cardRepository) Create(card *models.Card) error {
	return r.db.Create(card).Error
}

// ExistsByYGOProID checks if a Card with the given YGOProDeck ID already exists in the database.
func (r *cardRepository) ExistsByYGOProID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.Card{}).Where("card_ygo_id = ?", id).Count(&count).Error
	return count > 0, err
}
