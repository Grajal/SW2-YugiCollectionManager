package models

type Card struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	Desc      string
	FrameType string // Ex: "xyx", "spell", "trap"
	Type      string

	MonsterCard *MonsterCard   `gorm:"foreignKey:CardID"`
	SpellCard   *SpellTrapCard `gorm:"foreignKey:CardID"`
}
