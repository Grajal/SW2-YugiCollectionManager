package models

import "gorm.io/gorm"

type Deck struct {
	gorm.Model
	UserID      uint   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description string
	User        User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE, OnDelete:CASCADE"`

	Cards []DeckCard
}
