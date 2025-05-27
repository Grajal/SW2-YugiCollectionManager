package models

import "gorm.io/gorm"

type LinkMonsterCard struct {
	gorm.Model
	CardID      uint `gorm:"primaryKey;contraint: OnUpdate:CASCADE,OnDelete:CASCADE"`
	LinkValue   int
	LinkMarkers string `gorm:"type:jsonb"`
	Atk         int
	Level       int
	Attribute   string
	Race        string
}
