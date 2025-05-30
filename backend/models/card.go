package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	CardYGOID int `gorm:"unique;not null"`
	Name      string
	Desc      string
	FrameType string
	Type      string
	ImageURL  string

	MonsterCard         *MonsterCard
	SpellTrapCard       *SpellTrapCard
	LinkMonsterCard     *LinkMonsterCard
	PendulumMonsterCard *PendulumMonsterCard
}
