package models

type SpellTrapCard struct {
	CardID uint `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type   string
}
