package controller

import (
	"expensez/pkg/utils"
	v "expensez/pkg/validator"
	"expensez/src/models"
	"expensez/src/repository"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

var (
	ActivityCtrl Controller
)

type ActivityController struct {
	repo      repository.Repository
	validator v.CustomValidator
}

func ActivityLoad() error {
	ActivityCtrl = &ActivityController{repo: repository.Repo, validator: *v.Validator}
	return nil
}

type FullActivity struct {
	models.Activity
	FromAccount       string  `json:"from_account"`
	ToAccount         string  `json:"to_account"`
	Tag               string  `json:"tag"`
	SubTag            string  `json:"sub_tag"`
	TransactionType   string  `json:"transaction_type"`
	FpPreviousBalance float64 `json:"fp_previous_balance" gorm:"column:fp_previous_balance"`
	FpBalance         float64 `json:"fp_balance" gorm:"column:fp_balance"`
	TpPreviousBalance float64 `json:"tp_previous_balance" gorm:"column:tp_previous_balance"`
	TpBalance         float64 `json:"tp_balance" gorm:"column:tp_balance"`
}

func (t *ActivityController) Create(c echo.Context) error {
	modal := &models.Activity{}
	if err := c.Bind(modal); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	res, err := t.repo.Create(c, modal)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, res)
}

func (t *ActivityController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = t.repo.Delete(c, models.Activity{}, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Deleted")
}

func (t *ActivityController) Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	res, err := t.repo.Get(c, models.Activity{ID: uint(id)})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (t *ActivityController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	modl := &models.Activity{}
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

func (t *ActivityController) List(c echo.Context) error {
	tagID := c.QueryParam("tag_id")
	accountID := c.QueryParam("account_id")
	accountKey := c.QueryParam("account_key")
	transactionTypeID := c.QueryParam("transaction_type_id")
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	pageIndex := c.QueryParam("page_index")
	pageSize := c.QueryParam("page_size")
	remarks := c.QueryParam("remarks")

	conditions := []string{}
	if utils.IsValueNonZero(tagID) {
		condTag := fmt.Sprintf(" ( act.tag_id = %s or act.sub_tag_id = %s ) ", tagID, tagID)
		conditions = append(conditions, condTag)
	}
	if utils.IsValueNonZero(accountID) {
		condTag := fmt.Sprintf(" ( act.from_account_id = %s or act.to_account_id = %s ) ", accountID, accountID)
		if utils.IsValueNonZero(accountKey) {
			condTag = fmt.Sprintf(" ( ( act.from_account_id = %s and act.to_account_id = %s ) or ( act.from_account_id = %s and act.to_account_id = %s ) ) ", accountID, accountKey, accountKey, accountID)
		}
		conditions = append(conditions, condTag)
	}
	if utils.IsValueNonZero(transactionTypeID) {
		condTag := fmt.Sprintf(" act.transaction_type_id = %s ", transactionTypeID)
		conditions = append(conditions, condTag)
	}
	if utils.TrimSpace(startDate) != "" {
		condTag := fmt.Sprintf(" act.event_date >= %s ", startDate)
		conditions = append(conditions, condTag)
	}
	if utils.TrimSpace(endDate) != "" {
		condTag := fmt.Sprintf(" act.event_date <= %s ", endDate)
		conditions = append(conditions, condTag)
	}
	if utils.TrimSpace(remarks) != "" {
		condTag := " lower(act.remarks) like '%" + strings.ToLower(remarks) + "%' "
		conditions = append(conditions, condTag)
	}

	list := &[]FullActivity{}
	where := strings.Join(conditions[:], " AND ")
	limit := utils.StringToInt(pageSize)
	offset := (utils.StringToInt(pageIndex) - 1) * limit
	table := "`activities` as act"
	silect := `act.id, act.from_account_id, act.to_account_id, act.tag_id, act.amount, act.sub_tag_id, 
	DATE_FORMAT(act.event_date, "%Y-%m-%d") as event_date, act.remarks, act.transaction_type_id, act.created_at, act.updated_at,
	fa.name as from_account, ta.name as to_account, tg.name as tag, s_tg.name as sub_tag, 
	 transaction_types.name as transaction_type, fp.previous_balance as fp_previous_balance, 
	 fp.balance as fp_balance, tp.previous_balance as tp_previous_balance, tp.balance as tp_balance`
	joins := `LEFT JOIN tags tg ON tg.id = act.tag_id 
	LEFT JOIN tags s_tg ON s_tg.id = act.sub_tag_id 
	LEFT JOIN transaction_types ON transaction_types.id = act.transaction_type_id 
	LEFT JOIN passbooks as fp ON fp.activity_id = act.id and act.from_account_id = fp.account_id 
	LEFT JOIN passbooks as tp ON tp.activity_id = act.id and act.to_account_id = tp.account_id 
	LEFT JOIN accounts as fa ON fa.id = act.from_account_id 
	LEFT JOIN accounts as ta ON ta.id = act.to_account_id `
	orderBy := "`act`.`event_date` DESC, `act`.`id` DESC"
	err := t.repo.FetchWithFullQuery(c, list, table, silect, joins, where, "", orderBy, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var count, sum float64
	whereCond := ""
	if utils.TrimSpace(where) != "" {
		whereCond = fmt.Sprintf(" WHERE %s ", where)
	}
	t.repo.FetchRow("SELECT COUNT(act.id) as count FROM `activities` as act "+whereCond, &count)
	t.repo.FetchRow("SELECT SUM(act.amount) as sum FROM `activities` as act "+whereCond, &sum)

	return c.JSON(http.StatusOK, map[string]any{
		"list":  list,
		"total": count,
		"sum":   sum,
	})
}
