package controller

import (
	"expensez/src/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	repository "github.com/uzrnem/go/rdb"
)

var (
	PassbookCtrl Controller
)

type PassbookController struct {
	BaseController
}

func PassbookLoad() error {
	PassbookCtrl = &PassbookController{
		BaseController: BaseController{
			GetModel: func(id string) (any, error) {
				modal := &models.Passbook{}
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
				return &[]models.Passbook{}
			},
		},
	}
	return nil
}

func (t *PassbookController) Get(c echo.Context) error {
	list := &[]models.FullPassbook{}
	table := "passbooks p"
	silect := `ct.name as current_account, p.previous_balance, p.balance, tt.name as transaction_type, 
		CASE 
		WHEN ot.name is null THEN tg.name 
		WHEN tt.name = 'Credit' THEN CONCAT(tg.name, ' from ', ot.name) 
		ELSE CONCAT(tg.name, ' to ', ot.name) 
		END as comment, tg.name as tag_name, t.amount, 
		DATE_FORMAT(t.event_date, "%M %d, %Y") as event_date, 
		ot.name as opposite_account, t.remarks`
	joins := `left join transaction_types tt on p.transaction_type_id = tt.id 
		left join activities t on p.activity_id = t.id 
		left join tags tg on t.tag_id = tg.id 
		left join accounts ct on p.account_id = ct.id 
		left join accounts ot on (tt.name = 'Credit' and t.from_account_id = ot.id) or (tt.name = 'Debit' and t.to_account_id = ot.id)`
	orderBy := "t.event_date DESC, p.id DESC"
	where := fmt.Sprintf("p.account_id = %s", c.Param("id"))
	err := repository.Repo.Builder().Table(table).Select(silect).Join(joins).Where(where).Order(orderBy).Limit(15).Exec(list)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}
