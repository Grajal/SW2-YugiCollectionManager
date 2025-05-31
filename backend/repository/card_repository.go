package repository

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"gorm.io/gorm"
)

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

func (r *cardRepository) GetByID(id uint) (*models.Card, error) {
	var card models.Card
	err := r.db.Preload("MonsterCard").
		Preload("SpellTrapCard").
		Preload("LinkMonsterCard").
		Preload("PendulumMonsterCard").
		First(&card, "id = ?", id).Error

	return &card, err
}

func (r *cardRepository) GetByYGOProID(id int) (*models.Card, error) {
	var card models.Card
	err := r.db.Preload("MonsterCard").
		Preload("SpellTrapCard").
		Preload("LinkMonsterCard").
		Preload("PendulumMonsterCard").
		First(&card, "card_ygo_id = ?", id).Error
	return &card, err
}

func (r *cardRepository) GetByName(name string) (*models.Card, error) {
	var card models.Card
	err := r.db.Preload("MonsterCard").
		Preload("SpellTrapCard").
		Preload("LinkMonsterCard").
		Preload("PendulumMonsterCard").
		First(&card, "name ILIKE ?", name).Error
	return &card, err
}

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

func (r *cardRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.Card{}).Count(&count).Error
	return count, err
}

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

func (r *cardRepository) Create(card *models.Card) error {
	return r.db.Create(card).Error
}

func (r *cardRepository) ExistsByYGOProID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.Card{}).Where("card_ygo_id = ?", id).Count(&count).Error
	return count > 0, err
}
