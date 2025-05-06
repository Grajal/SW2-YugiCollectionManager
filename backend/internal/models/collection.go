package models

type Collection struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"not null"`
	Cards  []Card `gorm:"many2many:collection_cards;"`
}
