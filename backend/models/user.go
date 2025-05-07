package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string

	Collections []Collection `gorm:"foreignKey:UserID"`
	Decks       []Deck       `gorm:"foreignKey:UserID"`
}
