package models

type Card struct {
	ID        uint `gorm:"primaryKey"`
	CardYGOID int  `gorm:"unique;not null"`
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
