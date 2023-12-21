package controller

import (
	"expensez/src/models"
	"strconv"
)

var (
	TransactionTypeCtrl Controller
)

type TransactionTypeController struct {
	BaseController
}

func TransactionTypeLoad() error {
	TransactionTypeCtrl = &TransactionTypeController{
		BaseController: BaseController{
			GetModel: func(id string) (any, error) {
				modal := &models.TransactionType{}
				if id != "" {
					ID, err := strconv.Atoi(id)
					if err != nil {
						return nil, err
					}
					modal.ID = uint(ID)
				}
				return modal, nil
			},
			GetList: func() any {
				return &[]models.TransactionType{}
			},
		},
	}
	return nil
}
