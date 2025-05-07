// Package models contains data structures representing system entities
package models

// User represents a system user
// Contains basic user information and its relationships with collections and decks
type User struct {
	ID       uint   `gorm:"primaryKey"`      // Unique identifier
	Username string `gorm:"unique;not null"` // Unique username
	Email    string `gorm:"unique;not null"` // Unique email
	Password string // Hashed password

	Collections []Collection `gorm:"foreignKey:UserID"` // User's card collections
	Decks       []Deck       `gorm:"foreignKey:UserID"` // User's card decks
}
