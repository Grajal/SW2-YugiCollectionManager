package models

// UserCard represents the relationship between a user and a card in the collection.
// It includes fields for the unique ID of the relationship, the user's ID, the card's ID,
// and the quantity of the card the user owns.
// The struct also establishes a foreign key relationship with the Card model,
// ensuring that updates or deletions to a card are cascaded to the UserCard table.
type UserCard struct {
	UserID   uint `gorm:"primaryKey"`
	CardID   uint `gorm:"primaryKey"`
	Quantity int

	Card Card `gorm:"foreignKey:CardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
