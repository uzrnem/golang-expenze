package controller

import (
	v "expensez/pkg/validator"
	"expensez/src/models"
	"expensez/src/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

var (
	StatementCtrl Controller
)

type StatementController struct {
	repo      repository.Repository
	validator v.CustomValidator
}

func StatementLoad() error {
	StatementCtrl = &StatementController{repo: repository.Repo, validator: *v.Validator}
	return nil
}

type FullStatement struct {
	models.Statement
	Name string `json:"name" gorm:"column:name"`
}

func (t *StatementController) Create(c echo.Context) error {
	modal := &models.Statement{}
	if err := c.Bind(modal); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err := t.repo.Create(c, modal)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, modal)
}

func (t *StatementController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = t.repo.Delete(c, models.Statement{}, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Deleted")
}

func (t *StatementController) Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	res, err := t.repo.Get(c, models.Statement{ID: uint(id)})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (t *StatementController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	modl := &models.Statement{}
	if err := c.Bind(modl); err != nil {
		return err
	}
	modl.ID = uint(id)
	err = t.repo.Update(c, modl)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, modl)
}

func (t *StatementController) List(c echo.Context) error {
	list := &[]FullStatement{}
	table := "statements s"
	silect := `s.id, s.account_id, acc.name, s.amount, DATE_FORMAT(s.event_date, "%Y-%m-%d") as event_date, s.remarks`
	orderBy := "s.event_date desc"
	join := "LEFT JOIN accounts acc on s.account_id = acc.id"
	err := t.repo.FetchWithFullQuery(c, list, table, silect, join, "", "", orderBy, 100, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}
