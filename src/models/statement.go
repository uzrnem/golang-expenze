package models

import "time"

type Statement struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AccountID int64     `json:"account_id" gorm:"column:account_id"`
	Amount    float64   `json:"amount" gorm:"column:amount"`
	EventDate string    `json:"event_date" gorm:"column:event_date"`
	Remarks   string    `json:"remarks" gorm:"column:remarks"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:create"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at"`
}

type FullStatement struct {
	Statement
	Name string `json:"name" gorm:"column:name"`
}
