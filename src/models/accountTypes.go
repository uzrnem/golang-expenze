package models

type AccountType struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}
