package controller

import (
	"expensez/src/models"
	"strconv"
)

var (
	SubscriptionCtrl Controller
)

type SubscriptionController struct {
	BaseController
}

func SubscriptionLoad() error {
	SubscriptionCtrl = &SubscriptionController{
		BaseController: BaseController{
			GetModel: func(id string) (any, error) {
				modal := &models.Subscription{}
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
				return &[]models.Subscription{}
			},
		},
	}
	return nil
}
