package models

type User struct {
	ID       uint `gorm:"primaryKey"`
	Username string
	Password string

	Collections []Collection `gorm:"foreignKey:UserID"`
	Decks       []Deck       `gorm:"foreignKey:UserID"`
}
