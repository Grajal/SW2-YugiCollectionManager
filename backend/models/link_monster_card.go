package models

type LinkMonsterCard struct {
	CardID      uint `gorm:"primaryKey"`
	LinkValue   int
	LinkMarkers []string `gorm:"type:jsonb"`
	MonsterBase

	Card Card `gorm:"foreignKey:CardID;references:ID;contraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
