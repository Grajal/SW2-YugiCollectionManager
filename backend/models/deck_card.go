package models

import "gorm.io/gorm"

type DeckCard struct {
	gorm.Model
	DeckID   uint
	CardID   uint
	Quantity int    `gorm:"not null"`
	Zone     string `gorm:"type:varchar(10);not null"`

	Card Card `gorm:"foreignKey:CardID"`
}
