package models

import "time"

type Passbook struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	AccountID int64   `json:"account_id,string" gorm:"column:account_id"`
	ActivityID int64   `json:"activity_id,string" gorm:"column:activity_id"`
	TransactionTypeID int64  `json:"transaction_type_id,string" gorm:"column:transaction_type_id"`
	PreviousBalance          float64   `json:"previous_balance,string" gorm:"column:previous_balance"`
	Balance          float64   `json:"balance,string" gorm:"column:balance"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:create"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at"`
}