package models

import "time"

type Account struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name" gorm:"column:name;not null"`
	AccountTypeID     int64     `json:"account_type_id" gorm:"column:account_type_id"`
	Amount            float64   `json:"amount" gorm:"column:amount"`
	IsFrequent        bool      `json:"is_frequent" gorm:"column:is_frequent"`
	IsSnapshotDisable bool      `json:"is_snapshot_disable" gorm:"column:is_snapshot_disable"`
	IsClosed          bool      `json:"is_closed" gorm:"column:is_closed"`
	CreatedAt         time.Time `json:"created_at" gorm:"<-:create"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at"`
}

type FullAccount struct {
	Account
	Type string `json:"type" gorm:"column:type"`
}
