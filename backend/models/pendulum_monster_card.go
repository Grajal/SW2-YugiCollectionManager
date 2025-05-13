package models

type PendulumMonsterCard struct {
	CardID uint `gorm:"primaryKey"`
	Scale  int
	MonsterBase

	Card Card `gorm:"foreignKey:CardID;references:ID;contraint: OnUpdate:CASCADE,OnDelete:CASCADE"`
}
