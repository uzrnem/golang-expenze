package controller

import (
	pkgErr "expensez/pkg/errors"

	"github.com/labstack/echo"
)

var (
	ErrorMap = map[string]error{
		"Name_required":  pkgErr.New(pkgErr.VALIDATION_ERROR, "Name is required"),
		"Email_required": pkgErr.New(pkgErr.VALIDATION_ERROR, "Email is required"),
		"Email_email":    pkgErr.New(pkgErr.VALIDATION_ERROR, "Email is invalid"),
	}
)

type Controller interface {
	Create(c echo.Context) error
	Delete(c echo.Context) error
	Get(c echo.Context) error
	List(c echo.Context) error
	Update(c echo.Context) error
}

type ExtendedController interface {
	FindAccountByType(c echo.Context) error
	GetTagsByTranscationHits(c echo.Context) error
	BalanceSheet(c echo.Context) error
}