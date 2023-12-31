package models

type AccountType struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"column:name;not null"`
}

type FullAccountType struct {
	AccountType
	Amount string `json:"amount" gorm:"column:amount"`
}
