package models

type Deck struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	Name        string `gorm:"not null"`
	Description string

	DeckCards []DeckCard `gorm:"foreignKey:DeckID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
