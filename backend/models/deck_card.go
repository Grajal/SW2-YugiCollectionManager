package models

type DeckCard struct {
	DeckID   uint   `gorm:"primaryKey"`
	CardID   uint   `gorm:"primaryKey"`
	Quantity int    `gorm:"not null"`
	Zone     string `gorm:"type:varchar(10);not null"`

	Card Card `gorm:"foreignKey:CardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Deck Deck `gorm:"foreignKey:DeckID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
