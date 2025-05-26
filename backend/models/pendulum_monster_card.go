package models

type PendulumMonsterCard struct {
	CardID    uint `gorm:"primaryKey;contraint: OnUpdate:CASCADE,OnDelete:CASCADE"`
	Scale     int
	Atk       int
	Def       int
	Level     int
	Attribute string
	Race      string
}
