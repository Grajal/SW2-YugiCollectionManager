package models

type LinkMonsterCard struct {
	CardID      uint `gorm:"primaryKey"`
	LinkValue   int
	LinkMarkers []string `gorm:"type:jsonb"`
	Atk         int
	Level       int
	Attribute   string
	Race        string

	Card Card `gorm:"foreignKey:CardID;references:ID;contraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
