package repository

import "github.com/labstack/echo"

type Repository interface {
	Create(c echo.Context, modal any) (any, error)
	Delete(c echo.Context, modal any, id int) error
	Get(c echo.Context, modal any) (any, error)
	Update(c echo.Context, modal any) (any, error)
	List(c echo.Context, modal any) (any, error)
	FetchRow(sql string, res any)
	FetchWithQuery(c echo.Context, modal any, where, order string) error
	FetchWithFullQuery(c echo.Context, list any, table, silect, joins, where, groupBy, orderBy string, limit, offset int) error
}