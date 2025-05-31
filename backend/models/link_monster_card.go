package models

type LinkMonsterCard struct {
	CardID      uint `gorm:"primaryKey;contraint: OnUpdate:CASCADE,OnDelete:CASCADE"`
	LinkValue   int
	LinkMarkers string `gorm:"type:jsonb"`
	Atk         int
	Level       int
	Attribute   string
	Race        string
}
