package controller

import (
	"encoding/json"
	"expensez/src/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	repository "github.com/uzrnem/go/rdb"
	v "github.com/uzrnem/go/validator"
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
	modal := &models.Passbook{}
	if err := c.Bind(modal); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err := t.repo.Create(c.Request().Context(), modal)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(modal)
}

func (t *PassbookController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = t.repo.Delete(c.Request().Context(), models.Passbook{}, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Deleted")
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
	err := t.repo.Builder().Table(table).Select(silect).Join(joins).Where(where).Order(orderBy).Limit(15).Exec(list)
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
	err = t.repo.Update(c.Request().Context(), modl)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, modl)
}

func (t *PassbookController) List(c echo.Context) error {
	list := &[]models.Passbook{}
	res, err := t.repo.List(c.Request().Context(), list)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}
