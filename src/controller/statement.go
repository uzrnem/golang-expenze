package controller

import (
	"encoding/json"
	"expensez/src/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	repository "github.com/uzrnem/go/rdb"
	"github.com/uzrnem/go/utils"
	v "github.com/uzrnem/go/validator"
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

func (t *StatementController) Create(c echo.Context) error {
	modal := &models.Statement{}
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

func (t *StatementController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = t.repo.Delete(c.Request().Context(), models.Statement{}, id)
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
	res, err := t.repo.Get(c.Request().Context(), models.Statement{ID: uint(id)})
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
	err = t.repo.Update(c.Request().Context(), modl)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, modl)
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
	err := t.repo.Builder().Table(table).Select(silect).Join(join).Order(orderBy).Limit(limit).Offset(offset).Exec(list)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var count float64
	err = t.repo.Builder().Table(table).Select("COUNT(s.id) as count").Exec(&count)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"list":  list,
		"total": count,
	})
}
