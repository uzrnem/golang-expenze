package controller

import (
	"expensez/src/models"
	"strconv"
)

var (
	AccountTypeCtrl Controller
)

type AccountTypeController struct {
	BaseController
}

func AccountTypeLoad() error {
	AccountTypeCtrl = &AccountTypeController{
		BaseController: BaseController{
			GetModel: func(id string) (any, error) {
				modal := &models.AccountType{}
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
				return &[]models.AccountType{}
			},
		},
	}
	return nil
}
