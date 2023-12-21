package controller

import (
	"expensez/src/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	repository "github.com/uzrnem/go/rdb"
	"github.com/uzrnem/go/utils"
)

var (
	StockCtrl Controller
)

type StockController struct {
	BaseController
}

func StockLoad() error {
	StockCtrl = &StockController{
		BaseController: BaseController{
			GetModel: func(id string) (any, error) {
				modal := &models.Stock{}
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
				return &[]models.Stock{}
			},
		},
	}
	return nil
}

func (t *StockController) List(c echo.Context) error {
	pageIndex := c.QueryParam("page_index")
	pageSize := c.QueryParam("page_size")
	status := c.QueryParam("status")

	limit := utils.StringToInt(pageSize)
	offset := (utils.StringToInt(pageIndex) - 1) * limit
	where := "status = '" + status + "'"
	orderBy := "sell_date desc, buy_date desc"

	list := &[]models.Stock{}
	err := repository.Repo.Builder().Where(where).Order(orderBy).Limit(limit).Offset(offset).Exec(list)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var count float64
	err = repository.Repo.Builder().Table("stocks").Select("COUNT(id) as count").Where(where).Exec(&count)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"list":  list,
		"total": count,
	})
}
