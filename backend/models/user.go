// Package models contains data structures representing system entities
package models

import "gorm.io/gorm"

// User represents a system user
// Contains basic user information and its relationships with collections and decks
type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string // Hashed password

	Collections []Collection `gorm:"foreignKey:UserID"` // User's card collections
	Decks       []Deck       `gorm:"foreignKey:UserID"` // User's card decks
}
