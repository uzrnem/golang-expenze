package models

type TransactionType struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"column:name;not null"`
}
