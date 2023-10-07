package repository

import (
	"expensez/config/mysql"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

var (
	Repo Repository
)

type RepoStrct struct {
	MysqlDB *gorm.DB
}

func Load() error {
	Repo = &RepoStrct{MysqlDB: mysql.MysqlDB}
	return nil
}

// Create is a function that creates a new record in the repository.
//
// It takes in a `c` of type `echo.Context` and `modal` of type `any` as parameters.
// It returns a value of type `any` and an `error`.
func (r *RepoStrct) Create(c echo.Context, modal any) (any, error) {
	err := r.MysqlDB.Create(modal).Error
	if err != nil {
		return nil, err
	}
	return modal, nil
}

func (r *RepoStrct) Delete(c echo.Context, modal any, id int) error {
	err := r.MysqlDB.Delete(modal, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoStrct) Get(c echo.Context, modal any) (any, error) {
	err := r.MysqlDB.First(modal).Error
	if err != nil {
		return nil, err
	}
	return modal, nil
}

func (r *RepoStrct) Update(c echo.Context, modal any) (any, error) {
	err := r.MysqlDB.Save(modal).Error
	if err != nil {
		return nil, err
	}
	return modal, nil
}

func (r *RepoStrct) List(c echo.Context, modal any) (any, error) {
	err := r.MysqlDB.Find(modal).Error
	if err != nil {
		return nil, err
	}
	return modal, nil
}
