package controller

import (
	"encoding/json"
	"expensez/src/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	repository "github.com/uzrnem/go/rdb"
	v "github.com/uzrnem/go/validator"
)

var (
	AccountCtrl Controller
)

type AccountController struct {
	repo      repository.Repository
	validator v.CustomValidator
}

func AccountLoad() error {
	AccountCtrl = &AccountController{repo: repository.Repo, validator: *v.Validator}
	return nil
}

func (t *AccountController) Create(c echo.Context) error {
	modal := &models.Account{}
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

func (t *AccountController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = t.repo.Delete(c.Request().Context(), models.Account{}, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Deleted")
}

func (t *AccountController) Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	res, err := t.repo.Get(c.Request().Context(), models.Account{ID: uint(id)})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (t *AccountController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	modl := &models.Account{}
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

func (t *AccountController) List(c echo.Context) error {
	taip := c.QueryParam("type")
	list := &[]models.Account{}
	where := "is_closed = 0"
	order := "name ASC"
	if taip == "all_accounts" {
		where = ""
		order = "is_closed ASC, is_frequent DESC, is_snapshot_disable, account_type_id ASC, name ASC"
	} else if taip == "frequent" {
		where = "is_closed = 0 AND is_frequent = 1"
		order = "name ASC"
	}
	err := t.repo.Builder().Where(where).Order(order).Exec(list)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}
