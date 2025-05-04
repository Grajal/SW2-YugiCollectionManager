package models

type Deck struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"not null"`
	Name   string `gorm:"not null"`
	Cards  []Card `gorm:"many2many:deck_cards;"`
}
