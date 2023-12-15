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
	StockCtrl Controller
)

type StockController struct {
	repo      repository.Repository
	validator v.CustomValidator
}

func StockLoad() error {
	StockCtrl = &StockController{repo: repository.Repo, validator: *v.Validator}
	return nil
}

func (t *StockController) Create(c echo.Context) error {
	modal := &models.Stock{}
	if err := c.Bind(modal); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err := t.repo.Create(c.Request().Context(), modal)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(modal)
}

func (t *StockController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = t.repo.Delete(c.Request().Context(), models.Stock{}, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Deleted")
}

func (t *StockController) Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	res, err := t.repo.Get(c.Request().Context(), models.Stock{ID: uint(id)})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (t *StockController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	modl := &models.Stock{}
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

func (t *StockController) List(c echo.Context) error {
	pageIndex := c.QueryParam("page_index")
	pageSize := c.QueryParam("page_size")
	status := c.QueryParam("status")

	limit := utils.StringToInt(pageSize)
	offset := (utils.StringToInt(pageIndex) - 1) * limit
	where := "status = '" + status + "'"
	orderBy := "sell_date desc, buy_date desc"

	list := &[]models.Stock{}
	err := t.repo.Builder().Where(where).Order(orderBy).Limit(limit).Offset(offset).Exec(list)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var count float64
	err = t.repo.Builder().Table("stocks").Select("COUNT(id) as count").Where(where).Exec(&count)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"list":  list,
		"total": count,
	})
}
