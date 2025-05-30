package models

import "gorm.io/gorm"

type DeckCard struct {
	gorm.Model
	DeckID      uint
	CardID      uint
	Quantity    int `gorm:"not null"`
	IsExtraDeck bool

	Card Card `gorm:"foreignKey:CardID"`
}
