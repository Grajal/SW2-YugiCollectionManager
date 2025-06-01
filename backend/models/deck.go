package models

type Deck struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null;index"`
	Name        string `gorm:"not null"`
	Description string

	DeckCards []DeckCard `gorm:"foreignKey:DeckID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User      User       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
