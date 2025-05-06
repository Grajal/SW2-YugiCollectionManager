package models

type SpellTrapCard struct {
	CardID uint `gorm:"primaryKey"`
	Type   string

	Card Card `gorm:"foreignKey:CardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
