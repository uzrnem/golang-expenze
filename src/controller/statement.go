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
	StatementCtrl Controller
)

type StatementController struct {
	BaseController
}

func StatementLoad() error {
	StatementCtrl = &StatementController{
		BaseController: BaseController{
			GetModel: func(id string) (any, error) {
				modal := &models.Statement{}
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
				return &[]models.Statement{}
			},
		},
	}
	return nil
}

func (t *StatementController) List(c echo.Context) error {
	pageIndex := c.QueryParam("page_index")
	pageSize := c.QueryParam("page_size")
	limit := utils.StringToInt(pageSize)
	offset := (utils.StringToInt(pageIndex) - 1) * limit
	list := &[]models.FullStatement{}
	table := "statements s"
	silect := `s.id, s.account_id, acc.name, s.amount, DATE_FORMAT(s.event_date, "%Y-%m-%d") as event_date, s.remarks`
	orderBy := "s.event_date desc"
	join := "LEFT JOIN accounts acc on s.account_id = acc.id"
	err := repository.Repo.Builder().Table(table).Select(silect).Join(join).Order(orderBy).Limit(limit).Offset(offset).Exec(list)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var count float64
	err = repository.Repo.Builder().Table(table).Select("COUNT(s.id) as count").Exec(&count)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"list":  list,
		"total": count,
	})
}
