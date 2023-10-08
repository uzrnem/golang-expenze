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

func (t *StatementController) Create(c echo.Context) error {
	modal := models.Statement{}
	if errs := t.validator.Validate(modal); errs != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errs)
	}

	res, err := t.repo.Create(c, modal)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, res)
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
	res, err := t.repo.Update(c, modl)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (t *StatementController) List(c echo.Context) error {
	list := &[]models.Statement{}
	res, err := t.repo.List(c, list)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}
