package models

import "gorm.io/gorm"

// UserCard represents the relationship between a user and a card in the collection.
// It includes fields for the unique ID of the relationship, the user's ID, the card's ID,
// and the quantity of the card the user owns.
// The struct also establishes a foreign key relationship with the Card model,
// ensuring that updates or deletions to a card are cascaded to the UserCard table.
type UserCard struct {
	gorm.Model
	UserID   uint
	CardID   uint
	Quantity int

	Card Card `gorm:"foreignKey:CardID"`
}
