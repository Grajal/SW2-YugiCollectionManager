package models

import "gorm.io/gorm"

type MonsterCard struct {
	gorm.Model
	CardID    uint `gorm:"primaryKey;contraint: OnUpdate:CASCADE,OnDelete:CASCADE"`
	Atk       int
	Def       int
	Level     int
	Attribute string
	Race      string
}
