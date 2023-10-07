package models

type TransactionType struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}
