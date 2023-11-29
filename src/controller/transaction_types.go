package controller

import (
	v "expensez/pkg/validator"
	"expensez/src/models"
	"expensez/src/repository"
	"net/http"
	"strconv"
	"encoding/json"

	"github.com/labstack/echo"
)

var (
	TransactionTypeCtrl Controller
)

type TransactionTypeController struct {
	repo      repository.Repository
	validator v.CustomValidator
}

func TransactionTypeLoad() error {
	TransactionTypeCtrl = &TransactionTypeController{repo: repository.Repo, validator: *v.Validator}
	return nil
}

func (t *TransactionTypeController) Create(c echo.Context) error {
	modal := &models.TransactionType{}
	if err := c.Bind(modal); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err := t.repo.Create(c, modal)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(modal)
}

func (t *TransactionTypeController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = t.repo.Delete(c, models.TransactionType{}, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Deleted")
}

func (t *TransactionTypeController) Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	res, err := t.repo.Get(c, models.TransactionType{ID: uint(id)})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (t *TransactionTypeController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	modl := &models.TransactionType{}
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

func (t *TransactionTypeController) List(c echo.Context) error {
	list := &[]models.TransactionType{}
	res, err := t.repo.List(c, list)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}
