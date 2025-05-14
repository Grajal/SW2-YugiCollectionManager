package models

import "gorm.io/gorm"

type DeckCard struct {
	gorm.Model
	DeckID   uint `gorm:"not null"`
	CardID   uint `gorm:"not null"`
	Quantity int  `gorm:"not null"`

	IsExtraDeck bool

	Deck Deck `gorm:"foreignKey:DeckID;references:ID;constraint:OnUpdate:CASCADE, OnDelete:CASCADE"`
	Card Card `gorm:"foreignKey:CardID;references:ID;constraint:OnUpdate:CASCADE, OnDelete:CASCADE"`
}
