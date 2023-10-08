package controller

import (
	v "expensez/pkg/validator"
	"expensez/src/models"
	"expensez/src/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

var (
	PassbookCtrl Controller
)

type PassbookController struct {
	repo      repository.Repository
	validator v.CustomValidator
}

func PassbookLoad() error {
	PassbookCtrl = &PassbookController{repo: repository.Repo, validator: *v.Validator}
	return nil
}

func (t *PassbookController) Create(c echo.Context) error {
	modal := models.Passbook{}
	if errs := t.validator.Validate(modal); errs != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errs)
	}

	res, err := t.repo.Create(c, modal)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, res)
}

func (t *PassbookController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = t.repo.Delete(c, models.Passbook{}, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Deleted")
}

type FullPassbook struct {
	models.Passbook
	CurrentAccount  string  `json:"current_account"`
	TransactionType string  `json:"transaction_type"`
	Comment         string  `json:"comment"`
	TagName         string  `json:"tag_name"`
	Amount          float64 `json:"amount,string"`
	EventDate       string  `json:"event_date"`
	OppositeAccount string  `json:"opposite_account"`
	Remarks         string  `json:"remarks"`
}

func (t *PassbookController) Get(c echo.Context) error {
	list := &[]FullPassbook{}
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
	err := t.repo.FetchWithFullQuery(c, list, table, silect, joins, where, "", orderBy, 15, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}

func (t *PassbookController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	modl := &models.Passbook{}
	if err := c.Bind(modl); err != nil {
		return err
	}
	modl.ID = uint(id)
	res, err := t.repo.Update(c, modl)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (t *PassbookController) List(c echo.Context) error {
	list := &[]models.Passbook{}
	res, err := t.repo.List(c, list)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}
