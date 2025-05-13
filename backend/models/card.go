package models

type Card struct {
	ID        uint `gorm:"primaryKey"`
	CardYGOID int  `gorm:"unique;not null"`
	Name      string
	Desc      string
	FrameType string
	Type      string
}
