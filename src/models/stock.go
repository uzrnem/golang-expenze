package models

import (
	"errors"
	"time"
	"gorm.io/gorm"
)

type StatusType string

const (
	Open  StatusType = "open"
	Close StatusType = "close"
)

type Stock struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Code             string    `json:"code" gorm:"column:code;not null"`
	Description      string    `json:"description" gorm:"column:description;not null"`
	Quantity         uint      `json:"quantity" gorm:"column:quantity"`
	CurrentQuantity  uint      `json:"current_quantity" gorm:"column:current_quantity"`
	BuyPrice         float64   `json:"buy_price" gorm:"column:buy_price"`
	BuyCharges       float64   `json:"buy_charges" gorm:"column:buy_charges"`
	BuyDate          string    `json:"buy_date" gorm:"column:buy_date"`
	Target           float64   `json:"target" gorm:"column:target"`
	SellPrice        *float64   `json:"sell_price" gorm:"column:sell_price"`
	SellCharges      *float64   `json:"sell_charges" gorm:"column:sell_charges"`
	SellDate         *string    `json:"sell_date" gorm:"column:sell_date"`
	Status           StatusType  `json:"status" gorm:"column:status;type:enum('open', 'close')"`
	ParentID         *int      `json:"parent_id" gorm:"column:parent_id"`
	CreatedAt        time.Time `json:"created_at" gorm:"<-:create"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at"`
}

func (s *Stock) BeforeCreate(tx *gorm.DB) (err error) {
	if s.Status == "" {
		s.Status = "open"
	}
	return nil
}

func (s *Stock) AfterFind(tx *gorm.DB) (err error) {
	if len(s.BuyDate) > 10 {
		s.BuyDate = string(s.BuyDate[:10])
	}
	if s.SellDate != nil && len(*s.SellDate) > 10 {
		*s.SellDate = string((*s.SellDate)[:10])
	}
	return nil
}

func (s *Stock) BeforeUpdate(tx *gorm.DB) (err error) {
	dbStock := &Stock{ID: s.ID}
	err = tx.First(dbStock).Error
	if err != nil {
		return err
	}
	if dbStock.Status == "open" && s.Status == "close" && dbStock.CurrentQuantity < s.CurrentQuantity {
		return errors.New("not enough stock")
	}
	if dbStock.Status == "open" && s.Status == "close" && dbStock.CurrentQuantity > s.CurrentQuantity {
		parentId := int(dbStock.ID)
		if dbStock.ParentID != nil && *dbStock.ParentID > 0 {
			parentId = *dbStock.ParentID
		}
		pendingStock := &Stock{
			Status: "open",
			Code: dbStock.Code,
			Description: dbStock.Description,
			Quantity: dbStock.Quantity,
			CurrentQuantity: dbStock.CurrentQuantity - s.CurrentQuantity,
			BuyPrice: dbStock.BuyPrice,
			BuyCharges: dbStock.BuyCharges,
			BuyDate: dbStock.BuyDate,
			Target: dbStock.Target,
			ParentID: &parentId,
		}
		err = tx.Create(pendingStock).Error
		if err != nil {
			return err
		}
	}
	return nil
}
/* GORM Hooks
func (u *User) BeforeCreate(tx *gorm.DB) (err error)
func (u *User) AfterCreate(tx *gorm.DB) (err error)
func (u *User) BeforeUpdate(tx *gorm.DB) (err error)
func (u *User) AfterUpdate(tx *gorm.DB) (err error)
func (u *User) BeforeDelete(tx *gorm.DB) (err error)
func (u *User) AfterDelete(tx *gorm.DB) (err error)
func (s *Stock) AfterFind(tx *gorm.DB) (err error)
*/