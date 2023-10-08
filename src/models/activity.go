package models

import "time"

type Activity struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	FromAccountID int64   `json:"from_account_id" gorm:"column:from_account_id"`
	ToAccountID int64   `json:"to_account_id" gorm:"column:to_account_id"`
	TagID int64   `json:"tag_id" gorm:"column:tag_id"`
	Amount          float64   `json:"amount,string" gorm:"column:amount"`
	SubTagID int64   `json:"sub_tag_id" gorm:"column:sub_tag_id"`
	EventDate string `json:"event_date" gorm:"column:event_date"`
	Remarks string `json:"remarks" gorm:"column:remarks"`
	TransactionTypeID int64   `json:"transaction_type_id" gorm:"column:transaction_type_id"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:create"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at"`
}