package models

import "gorm.io/gorm"

type Deck struct {
	gorm.Model
	UserID      uint
	Name        string `gorm:"not null"`
	Description string

	DeckCards []DeckCard `gorm:"foreignKey:DeckID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
