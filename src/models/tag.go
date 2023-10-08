package models

type Tag struct {
	ID                uint   `json:"id" gorm:"primaryKey"`
	Name              string `json:"name" gorm:"column:name;not null"`
	TransactionTypeID int64  `json:"transaction_type_id,string" gorm:"column:transaction_type_id"`
	TagID             int64  `json:"tag_id" gorm:"column:tag_id"`
}