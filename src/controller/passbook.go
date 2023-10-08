package controller

import (
	"fmt"
	v "expensez/pkg/validator"
	"expensez/src/models"
	"expensez/src/repository"
	"net/http"
	"strconv"
	"github.com/labstack/echo"
	"expensez/pkg/utils"
	"strings"
)

var (
	PassbookCtrl Controller
)

type PassbookController struct {
	repo      repository.Repository
	validator v.CustomValidator
}

func PassbookLoad() error {
	PassbookCtrl = &PassbookController{repo: repository.Repo, validator: *v.Validator}
	return nil
}

func (t *PassbookController) Create(c echo.Context) error {
	modal := models.Passbook{}
	if errs := t.validator.Validate(modal); errs != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errs)
	}

	res, err := t.repo.Create(c, modal)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, res)
}

func (t *PassbookController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = t.repo.Delete(c, models.Passbook{}, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Deleted")
}

type FullPassbook struct {
	models.Passbook
	CurrentAccount string `json:"current_account"`
	TransactionType          string   `json:"transaction_type"`
	Comment string `json:"comment"`
	TagName string `json:"tag_name"`
	Amount          float64   `json:"amount,string"`
	EventDate string `json:"event_date"`
	OppositeAccount string `json:"opposite_account"`
	Remarks string `json:"remarks"`
}

func (t *PassbookController) Get(c echo.Context) error {
	list := &[]FullPassbook{}
	table := "passbooks p"
	silect := `ct.name as current_account, p.previous_balance, p.balance, tt.name as transaction_type, 
		CASE 
		WHEN ot.name is null THEN tg.name 
		WHEN tt.name = 'Credit' THEN CONCAT(tg.name, ' from ', ot.name) 
		ELSE CONCAT(tg.name, ' to ', ot.name) 
		END as comment, tg.name as tag_name, t.amount, 
		DATE_FORMAT(t.event_date, "%M %d, %Y") as event_date, 
		ot.name as opposite_account, t.remarks`
	joins := `left join transaction_types tt on p.transaction_type_id = tt.id 
		left join activities t on p.activity_id = t.id 
		left join tags tg on t.tag_id = tg.id 
		left join accounts ct on p.account_id = ct.id 
		left join accounts ot on (tt.name = 'Credit' and t.from_account_id = ot.id) or (tt.name = 'Debit' and t.to_account_id = ot.id)`
	orderBy := "t.event_date DESC, p.id DESC"
	where := fmt.Sprintf("p.account_id = %d", c.Param("id"))
	err := t.repo.FetchWithFullQuery(c, list, table, silect, joins, where, "", orderBy, 15, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}

func (t *PassbookController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	modl := &models.Passbook{}
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

func (t *PassbookController) List(c echo.Context) error {
	tagID := c.QueryParam("tag_id")
	accountID := c.QueryParam("account_id")
	accountKey := c.QueryParam("account_key")
	transactionTypeID := c.QueryParam("transaction_type_id")
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	pageIndex := c.QueryParam("page_index")
	pageSize := c.QueryParam("page_size")


	conditions := []string{}
	if utils.IsValueNonZero(tagID) {
		condTag := fmt.Sprintf(" ( act.tag_id = %s or act.sub_tag_id = %s ) ", tagID, tagID)
		conditions = append(conditions, condTag)
	}
	if utils.IsValueNonZero(accountID) {
		condTag := fmt.Sprintf(" ( act.from_account_id = %s or act.to_account_id = %s ) ", accountID, accountID)
		if utils.IsValueNonZero(accountKey) {
			condTag = fmt.Sprintf(" ( ( act.from_account_id = %s and act.to_account_id = %s ) or ( act.from_account_id = $s and act.to_account_id = %s ) ) ", accountID, accountKey, accountKey, accountID)
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

	list := &[]FullActivity{}
	where := strings.Join(conditions[:], " AND ")
	limit := utils.StringToInt(pageSize)
	offset := (utils.StringToInt(pageIndex) - 1 ) * limit
	table := "`activities` as act"
	silect := `act.*, fa.name as from_account, ta.name as to_account, tg.name as tag, s_tg.name as sub_tag, 
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
	var count, sum int
	whereCond := ""
	if utils.TrimSpace(endDate) != "" {
		whereCond = fmt.Sprintf(" WHERE %s ", where)
	}

	t.repo.FetchRow("SELECT COUNT(act.id) as count FROM `activities` as act " + whereCond, &count)
	t.repo.FetchRow("SELECT SUM(act.amount) as sum FROM `activities` as act " + whereCond, &sum)

	return c.JSON(http.StatusOK, map[string]any{
		"list": list,
		"total": count,
		"sum": sum,
	})
}