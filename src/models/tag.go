package models

type Tag struct {
	ID                uint   `json:"id" gorm:"primaryKey"`
	Name              string `json:"name" gorm:"not null"`
	TransactionTypeID uint   `json:"transaction_type_id" gorm:"column:transaction_type_id"`
	TagID             uint   `json:"tag_id" gorm:"column:tag_id"`
}
