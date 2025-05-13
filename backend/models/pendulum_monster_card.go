package models

type PendulumMonsterCard struct {
	CardID    uint `gorm:"primaryKey"`
	Scale     int
	Atk       int
	Def       int
	Level     int
	Attribute string
	Race      string

	Card Card `gorm:"foreignKey:CardID;references:ID;contraint: OnUpdate:CASCADE,OnDelete:CASCADE"`
}
