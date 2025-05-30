package models

import "gorm.io/gorm"

type SpellTrapCard struct {
	gorm.Model
	CardID uint `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type   string
}
