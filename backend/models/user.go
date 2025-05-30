// Package models contains data structures representing system entities
package models

import "gorm.io/gorm"

// User represents a system user
// Contains basic user information and its relationships with collections and decks
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string // Hashed password

	Collection []UserCard `gorm:"primaryKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Decks      []Deck     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
