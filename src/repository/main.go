package repository

import (
	//"fmt"
	"strings"
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


/*
	sql := r.MysqlDB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		db := tx.Table(table)
		return db.Find(list)
	  })
	fmt.Println("SQL : ", sql)
	rows, err := db.Having("sum(amount) > ?", 100).Rows()
&{
	DB:0x4000171920 
	TableExpr:0x40001719e0 
	Table:act Model:0x400008ed80 
	Unscoped:false 
	Dest:0x400008ed80 
	ReflectValue:<[]controller.FullActivity Value> 
	Clauses:map[
		FROM:{
			Name:FROM 
			BeforeExpression:<nil>
			AfterNameExpression:<nil> 
			AfterExpression:<nil> 
			Expression:{
				Tables:[] 
				Joins:[
					{
						Type: 
						Table:{
							Name: 
							Alias: 
							Raw:false
						}
						ON:{
							Exprs:[]
						}
						Using:[]
						Expression:{
							SQL:LEFT JOIN tags tg ON tg.id = act.tag_id 
	LEFT JOIN tags s_tg ON s_tg.id = act.sub_tag_id 
	LEFT JOIN transaction_types ON transaction_types.id = act.transaction_type_id 
	LEFT JOIN accounts as fa ON fa.id = act.from_account_id 
	LEFT JOIN accounts as ta ON ta.id = act.to_account_id Vars:[]
						}
					}
				]
			}
			Builder:<nil>
		}
		LIMIT:{
			Name: BeforeExpression:<nil> 
			AfterNameExpression:<nil> 
			AfterExpression:<nil> 
			Expression:{Limit:0x400000eb20 Offset:0}
			Builder:<nil>
		}
		ORDER BY:{
			Name:ORDER BY 
			BeforeExpression:<nil> 
			AfterNameExpression:<nil> 
			AfterExpression:<nil> 
			Expression:{
				Columns:[
					{
						Column:{
							Table: 
							Name:act.created_at DESC, act.event_date DESC, act.id DESC 
							Alias: 
							Raw:true
						}
						Desc:false
						Reorder:false
					}
				]
				Expression:<nil>
			}
			Builder:<nil>
		}
		SELECT:{
			Name:SELECT 
			BeforeExpression:<nil> 
			AfterNameExpression:<nil> 
			AfterExpression:<nil> 
			Expression:{
				Distinct:false 
				Columns:[
					{
						Table: 
						Name:act.*, 
	fa.name as from_account, ta.name as to_account, tg.name as tag, s_tg.name as sub_tag, 
	transaction_types.name as transaction_type 
						Alias:
						Raw:true
					}
				]
				Expression:<nil>
			}
			Builder:<nil>
		}
	]
	BuildClauses:[]
	Distinct:false
	Selects:[act.*, 
	fa.name as from_account, ta.name as to_account, tg.name as tag, s_tg.name as sub_tag, 
	transaction_types.name as transaction_type]
	Omits:[] 
	Joins:[] 
	Preloads:map[] 
	Settings:{mu:{state:0 sema:0} read:{_:[] _:{} v:<nil>} dirty:map[] misses:0} ConnPool:0x4000092a90 Schema:expensez/src/controller.FullActivity Context:context.Background RaiseErrorOnNotFound:false SkipHooks:false SQL:{addr:<nil> buf:[]} Vars:[] CurDestIndex:0 attrs:[] assigns:[] scopes:[]
}
*/