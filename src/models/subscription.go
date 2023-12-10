package models

import "time"

type PaidType string

const (
	Prepaid  PaidType = "prepaid"
	Postpaid PaidType = "postpaid"
)

type Subscription struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title" gorm:"column:title;not null"`
	Description  string    `json:"description" gorm:"column:description;not null"`
	Amount       float64   `json:"amount" gorm:"column:amount"`
	ActualAmount float64   `json:"actual_amount" gorm:"column:actual_amount"`
	StartDate    string    `json:"start_date" gorm:"column:start_date"`
	EndDate      string    `json:"end_date" gorm:"column:end_date"`
	DueDate      string    `json:"due_date" gorm:"column:due_date"`
	Type         PaidType  `json:"type" gorm:"column:type;type:enum('prepaid', 'postpaid')"`
	Status       bool      `json:"status" gorm:"column:status"`
	CreatedAt    time.Time `json:"created_at" gorm:"<-:create"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at"`
}
