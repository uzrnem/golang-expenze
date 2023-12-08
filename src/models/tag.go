package models

type Tag struct {
	ID                uint   `json:"id" gorm:"primaryKey"`
	Name              string `json:"name" gorm:"column:name;not null"`
	TransactionTypeID int64  `json:"transaction_type_id" gorm:"column:transaction_type_id"`
	TagID             *int   `json:"tag_id" gorm:"column:tag_id"`
}

type FullTag struct {
	Tag
	Parent   string `json:"parent"`
	Type     string `json:"type"`
	TagCount uint   `json:"tag_count"`
}
