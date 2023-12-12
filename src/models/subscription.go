package models

import (
	"time"
	"gorm.io/gorm"
)

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
	StartDate    *string   `json:"start_date" gorm:"column:start_date"`
	EndDate      string    `json:"end_date" gorm:"column:end_date"`
	DueDate      *string   `json:"due_date" gorm:"column:due_date"`
	Type         PaidType  `json:"type" gorm:"column:type;type:enum('prepaid', 'postpaid')"`
	Status       bool      `json:"status" gorm:"column:status"`
	CreatedAt    time.Time `json:"created_at" gorm:"<-:create"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at"`
}

func (s *Subscription) AfterFind(tx *gorm.DB) (err error) {
	if s.StartDate != nil && len(*s.StartDate) > 10 {
		*s.StartDate = string((*s.StartDate)[:10])
	}
	if len(s.EndDate) > 10 {
		s.EndDate = string(s.EndDate[:10])
	}
	if s.DueDate != nil && len(*s.DueDate) > 10 {
		*s.DueDate = string((*s.DueDate)[:10])
	}
	return nil
}
