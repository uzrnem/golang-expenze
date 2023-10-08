package repository

import (
	"expensez/config/mysql"
	"strings"

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
func (r *RepoStrct) Create(c echo.Context, modal any) error {
	err := r.MysqlDB.Create(modal).Error
	if err != nil {
		return err
	}
	return nil
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

func (r *RepoStrct) Update(c echo.Context, modal any) error {
	err := r.MysqlDB.Save(modal).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoStrct) List(c echo.Context, modal any) (any, error) {
	err := r.MysqlDB.Find(modal).Error
	if err != nil {
		return nil, err
	}
	return modal, nil
}

func (r *RepoStrct) FetchWithQuery(c echo.Context, modal any, where, order string) error {
	return r.MysqlDB.Debug().Model(modal).Where(where).Order(order).Find(modal).Error
}

func (r *RepoStrct) FetchWithFullQuery(c echo.Context, list any, table, silect, joins, where, groupBy, orderBy string, limit, offset int) error {
	db := r.MysqlDB.Debug().Table(table)
	if strings.TrimSpace(silect) != "" {
		db = db.Select(silect)
	}
	if strings.TrimSpace(joins) != "" {
		db = db.Joins(joins)
	}
	if strings.TrimSpace(where) != "" {
		db = db.Where(where)
	}
	if strings.TrimSpace(groupBy) != "" {
		db = db.Group(groupBy)
	}
	if strings.TrimSpace(orderBy) != "" {
		db = db.Order(orderBy)
	}
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	return db.Find(list).Error
}

func (r *RepoStrct) FetchRow(sql string, res any) {
	r.MysqlDB.Debug().Raw(sql).Scan(res)
}
