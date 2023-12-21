package controller

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	repository "github.com/uzrnem/go/rdb"
)

var (
	BaseCtrl Controller
)

type BaseController struct {
	GetModel func(id string) (any, error)
	GetList  func() any
}

func (t *BaseController) Create(c echo.Context) error {
	modal, err := t.GetModel("")
	if err := c.Bind(modal); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = repository.Repo.Create(c.Request().Context(), modal)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(modal)
}

func (t *BaseController) Delete(c echo.Context) error {
	modal, err := t.GetModel(c.Param("id"))
	err = repository.Repo.DeleteOne(c.Request().Context(), modal)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Deleted")
}

func (t *BaseController) Get(c echo.Context) error {
	modal, err := t.GetModel(c.Param("id"))
	res, err := repository.Repo.Get(c.Request().Context(), modal)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (t *BaseController) Update(c echo.Context) error {
	modal, err := t.GetModel(c.Param("id"))
	if err := c.Bind(modal); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = repository.Repo.Update(c.Request().Context(), modal)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, modal)
}

func (t *BaseController) List(c echo.Context) error {
	res, err := repository.Repo.List(c.Request().Context(), t.GetList())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}
