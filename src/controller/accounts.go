package controller

import (
	"expensez/src/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	repository "github.com/uzrnem/go/rdb"
)

var (
	AccountCtrl Controller
)

type AccountController struct {
	BaseController
}

func AccountLoad() error {
	AccountCtrl = &AccountController{
		BaseController: BaseController{
			GetModel: func(id string) (any, error) {
				modal := &models.Account{}
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
				return &[]models.Account{}
			},
		},
	}
	return nil
}

func (t *AccountController) List(c echo.Context) error {
	taip := c.QueryParam("type")
	list := &[]models.Account{}
	where := "is_closed = 0"
	order := "name ASC"
	if taip == "all_accounts" {
		where = ""
		order = "is_closed ASC, is_frequent DESC, is_snapshot_disable, account_type_id ASC, name ASC"
	} else if taip == "frequent" {
		where = "is_closed = 0 AND is_frequent = 1"
		order = "name ASC"
	}
	err := repository.Repo.Builder().Where(where).Order(order).Exec(list)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}
