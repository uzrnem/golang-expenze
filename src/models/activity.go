package models

import "time"

type Activity struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	FromAccountID     int64     `json:"from_account_id" gorm:"column:from_account_id"`
	ToAccountID       int64     `json:"to_account_id" gorm:"column:to_account_id"`
	TagID             int64     `json:"tag_id" gorm:"column:tag_id"`
	Amount            float64   `json:"amount" gorm:"column:amount"`
	SubTagID          int64     `json:"sub_tag_id" gorm:"column:sub_tag_id"`
	EventDate         string    `json:"event_date" gorm:"column:event_date"`
	Remarks           string    `json:"remarks" gorm:"column:remarks"`
	TransactionTypeID int64     `json:"transaction_type_id" gorm:"column:transaction_type_id"`
	CreatedAt         time.Time `json:"created_at" gorm:"<-:create"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at"`
}

type FullActivity struct {
	Activity
	FromAccount       string  `json:"from_account"`
	ToAccount         string  `json:"to_account"`
	Tag               string  `json:"tag"`
	SubTag            string  `json:"sub_tag"`
	TransactionType   string  `json:"transaction_type"`
	FpPreviousBalance float64 `json:"fp_previous_balance" gorm:"column:fp_previous_balance"`
	FpBalance         float64 `json:"fp_balance" gorm:"column:fp_balance"`
	TpPreviousBalance float64 `json:"tp_previous_balance" gorm:"column:tp_previous_balance"`
	TpBalance         float64 `json:"tp_balance" gorm:"column:tp_balance"`
}
