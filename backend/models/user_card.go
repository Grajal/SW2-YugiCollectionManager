package models

type UserCard struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint `gorm:"not null"`
	CardID   uint `gorm:"not null"`
	Quantity int

	Card Card `gorm:"foreignKey:CardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
