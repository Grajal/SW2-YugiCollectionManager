package models

type MonsterCard struct {
	CardID    uint `gorm:"primaryKey"`
	Atk       int
	Def       int
	Level     int
	Attribute string
	Race      string

	Card Card `gorm:"foreignKey:CardID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
