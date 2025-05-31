package models

type MonsterCard struct {
	CardID    uint `gorm:"primaryKey;contraint: OnUpdate:CASCADE,OnDelete:CASCADE"`
	Atk       int
	Def       int
	Level     int
	Attribute string
	Race      string
}
