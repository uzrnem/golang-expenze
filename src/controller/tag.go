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
	TagCtrl Controller
)

type TagController struct {
	BaseController
}

func TagLoad() error {
	TagCtrl = &TagController{
		BaseController: BaseController{
			GetModel: func(id string) (any, error) {
				modal := &models.Tag{}
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
				return &[]models.Tag{}
			},
		},
	}
	return nil
}

func (t *TagController) List(c echo.Context) error {
	parentTag := c.QueryParam("parentTag")
	list := &[]models.FullTag{}
	where := ""
	if utils.IsValueNonZero(parentTag) {
		where = parentTag + " IN (t.id, t.tag_id)"
	}
	table := "tags t"
	silect := `t.id, t.name, t.tag_id, p.name AS parent, t.transaction_type_id, 
	tt.name AS type, COUNT(DISTINCT(m.id)) AS tag_count`
	joins := `LEFT JOIN tags p ON t.tag_id = p.id 
	LEFT JOIN tags c ON t.id = c.tag_id 
	LEFT JOIN transaction_types tt ON t.transaction_type_id = tt.id 
	LEFT JOIN activities m ON t.id in (m.tag_id, m.sub_tag_id) AND m.event_date > DATE_SUB(now(), INTERVAL 6 MONTH)`
	groupBy := "t.id, t.name, t.tag_id, p.name, t.transaction_type_id, tt.name"
	orderBy := "COUNT(DISTINCT(c.id)) DESC, t.name ASC"

	err := repository.Repo.Builder().Select(silect).Table(table).Join(joins).Where(where).Group(groupBy).Order(orderBy).Exec(list)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}
