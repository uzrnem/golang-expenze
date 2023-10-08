package controller

import (
	"expensez/pkg/utils"
	v "expensez/pkg/validator"
	"expensez/src/models"
	"expensez/src/repository"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

var (
	ExtendedCtrl ExtendedController
)

type ExtController struct {
	repo      repository.Repository
	validator v.CustomValidator
}

func ExtendedLoad() error {
	ExtendedCtrl = &ExtController{repo: repository.Repo, validator: *v.Validator}
	return nil
}

func (t *ExtController) FindAccountByType(c echo.Context) error {
	account_type := c.Param("accountType")
	fmt.Println("account_type: " + account_type)
	list := &[]models.Account{}
	where := "account_type_id in (SELECT id FROM account_types where lower(name) = '" + strings.ToLower(account_type) + "')"
	order := "name ASC"
	err := t.repo.FetchWithQuery(c, list, where, order)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}

func (t *ExtController) GetTagsByTranscationHits(c echo.Context) error {
	fromAccountID := c.QueryParam("from_account_id")
	toAccountID := c.QueryParam("to_account_id")
	tagID := c.QueryParam("tag_id")
	transactionID := 1
	if !utils.IsValueNonZero(fromAccountID) {
		transactionID = 3 //income
	} else if !utils.IsValueNonZero(toAccountID) {
		transactionID = 2 //debit
	}
	mapp := map[string]any{
		"tags":       nil,
		"sub_tags":   nil,
		"tag_id":     utils.StringToInt(tagID),
		"sub_tag_id": 0,
		"remarks":    "",
	}
	if fromAccountID == "0" && toAccountID == "0" {
		mapp["tag_id"] = 0
		return c.JSON(http.StatusOK, mapp)
	}
	remarkQuery := " event_date > DATE_SUB(now(), INTERVAL 6 MONTH) "
	conditions := []string{remarkQuery}

	fromQuery := " from_account_id is null "
	if utils.IsValueNonZero(fromAccountID) {
		fromQuery = " from_account_id = " + fromAccountID + " "
	}
	conditions = append(conditions, fromQuery)

	toQuery := " to_account_id is null "
	if utils.IsValueNonZero(toAccountID) {
		toQuery = " to_account_id = " + toAccountID + " "
	}
	conditions = append(conditions, toQuery)

	//Transaction Id Condition
	transactionQuery := " transaction_type_id = " + fmt.Sprintf("%d", transactionID) + " "
	conditions = append(conditions, transactionQuery)

	tagList := &[]models.Tag{}
	where := transactionQuery + " AND tag_id IS NULL "
	err := t.repo.FetchWithQuery(c, tagList, where, "")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	mapp["tags"] = tagList

	if !utils.IsValueNonZero(tagID) {
		tagConditions := strings.Join(conditions[:], " AND ")
		actRes := &models.Activity{}
		table := "activities"
		silect := "tag_id"
		err := t.repo.FetchWithFullQuery(c, actRes, table, silect, "", tagConditions, "tag_id, sub_tag_id", "count(id) desc", 1, 0)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		tagID = fmt.Sprintf("%d", actRes.TagID)
		mapp["tag_id"] = actRes.TagID
	}
	if utils.IsValueNonZero(tagID) {
		conditions = append(conditions, "tag_id = "+tagID)
		subTagConditions := strings.Join(conditions[:], " AND ")
		actRes := &models.Activity{}
		table := "activities"
		silect := "sub_tag_id, GROUP_CONCAT(DISTINCT(remarks)) as remarks"
		err := t.repo.FetchWithFullQuery(c, actRes, table, silect, "", subTagConditions, "sub_tag_id", "count(id) desc", 1, 0)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		mapp["sub_tag_id"] = actRes.SubTagID
		mapp["remarks"] = actRes.Remarks

		tagList := &[]models.Tag{}
		err = t.repo.FetchWithQuery(c, tagList, "tag_id = "+tagID, "")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		mapp["sub_tags"] = tagList
	}
	return c.JSON(http.StatusOK, mapp)
}

func (t *ExtController) BalanceSheet(c echo.Context) error {
	actTypeRes := &[]FullAccountType{}
	table := "accounts a"
	silect := "t.name as name, SUM(a.amount) as amount "
	joins := "LEFT JOIN account_types t on a.account_type_id = t.id"
	where := "a.amount !=0 and a.is_snapshot_disable = 0 and a.is_closed != 1 "
	groupBy := `a.account_type_id order by t.name='Saving' desc, t.name='Credit' desc, t.name='Wallet' desc,
	t.name='Stocks Equity' desc, t.name='Loan' desc, t.name='Mutual Funds' desc, t.name='Deposit' desc`
	err := t.repo.FetchWithFullQuery(c, actTypeRes, table, silect, joins, where, groupBy, "", 0, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	actRes := &[]FullAccount{}
	silect = "a.name as name, t.name as type, a.amount as amount"
	orderBy := `t.name='Saving' desc, t.name='Credit' desc, t.name='Wallet' desc,
	t.name='Deposit' desc, t.name='Loan' desc, t.name='Stocks Equity', a.name`
	err = t.repo.FetchWithFullQuery(c, actRes, table, silect, joins, where, "", orderBy, 0, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"holding_balance": actTypeRes,
		"account_balance": actRes,
	})
}

type StatementModel struct {
	Year    int     `json:"year" gorm:"column:year"`
	Mon     string  `json:"mon" gorm:"column:mon"`
	Salary  float64 `json:"salary" gorm:"column:salary"`
	Count   int     `json:"count" gorm:"column:count"`
	Expense float64 `json:"expense" gorm:"column:expense"`
	Credit  string  `json:"credit" gorm:"column:credit"`
	Bill    float64 `json:"bill" gorm:"column:bill"`
	CCState string  `json:"r_cc" gorm:"column:r_cc"`
	Income  float64 `json:"income" gorm:"column:income"`
}

func (t *ExtController) Statement(c echo.Context) error {
	duration := c.Param("duration")

	date_condition := fmt.Sprintf(" event_date > MAKEDATE( YEAR( DATE_SUB(now(), INTERVAL %s YEAR) ), DAYOFYEAR( DATE_SUB(now(), INTERVAL DAYOFMONTH(now()) DAY ) ) + 1 ) ", duration)
	and_date_condition := ""
	where_date_condition := ""
	if utils.IsValueNonZero(duration) {
		and_date_condition = " and" + date_condition
		where_date_condition = " where" + date_condition
	}

	list := &[]StatementModel{}
	silect := `yrmn.year, yrmn.mon, SUM(salr.amount) AS salary, COUNT(salr.amount) AS count,
	SUM(exps.amount) AS expense, GROUP_CONCAT(stmn.cc) as credit, SUM(stmn.bill) as bill,
	GROUP_CONCAT(ccexp.exp) as r_cc, SUM(incm.amount) as income`
	table := fmt.Sprintf(`(
		SELECT DISTINCT YEAR(event_date) AS year, MONTHNAME(event_date) AS mon,
			EXTRACT(YEAR_MONTH From event_date) AS yearmonth
		FROM activities %s
		GROUP BY EXTRACT(YEAR_MONTH FROM event_date), YEAR(event_date), MONTHNAME(event_date)
	 ) yrmn`, where_date_condition)
	joins := fmt.Sprintf(`LEFT JOIN (
		SELECT SUM(amount) AS amount, EXTRACT(YEAR_MONTH FROM event_date) AS event_date
		FROM activities WHERE from_account_id is null AND transaction_type_id IN (
	         SELECT id FROM transaction_types WHERE name = 'Income') %s
	     GROUP BY EXTRACT(YEAR_MONTH FROM event_date)
	 ) incm ON incm.event_date = yrmn.yearmonth
	 LEFT JOIN (
		SELECT SUM(amount) AS amount, EXTRACT(YEAR_MONTH FROM event_date) AS event_date
		FROM activities WHERE sub_tag_id IN (SELECT id FROM tags WHERE name = 'Salary')  %s
	     GROUP BY EXTRACT(YEAR_MONTH FROM event_date)
	 ) salr ON salr.event_date = yrmn.yearmonth
	 LEFT JOIN (
		SELECT GROUP_CONCAT(DISTINCT(CONCAT(a.name, ':', s.amount, ':', IFNULL(s.remarks, '')))) AS cc,
		   EXTRACT(YEAR_MONTH FROM s.event_date) AS event_date, SUM(s.amount) as bill
		FROM statements s LEFT JOIN accounts a ON s.account_id = a.id %s
		GROUP BY EXTRACT(YEAR_MONTH FROM s.event_date)
	 ) stmn ON stmn.event_date = yrmn.yearmonth
	 LEFT JOIN (
		SELECT SUM(amount) AS amount, EXTRACT(YEAR_MONTH FROM event_date) AS event_date
		FROM activities WHERE transaction_type_id IN (
	         SELECT id FROM transaction_types WHERE name = 'Expense') %s
		GROUP BY EXTRACT(YEAR_MONTH FROM event_date)
	 ) exps ON exps.event_date = yrmn.yearmonth
	 LEFT JOIN (
	  Select GROUP_CONCAT(cc.exp) as exp, yearmonth from (
		SELECT concat(name, ':', sum(a.amount)) as exp, EXTRACT(YEAR_MONTH FROM event_date) as yearmonth
		FROM activities a
		LEFT JOIN accounts acc ON a.from_account_id = acc.id AND a.to_account_id IS NULL
		WHERE acc.account_type_id IN (SELECT id FROM account_types WHERE name = 'Credit') %s
		GROUP BY EXTRACT(YEAR_MONTH FROM event_date), acc.name
	  ) as cc GROUP BY yearmonth
	 ) ccexp on ccexp.yearmonth = yrmn.yearmonth `, and_date_condition, and_date_condition, where_date_condition, and_date_condition, and_date_condition)
	groupBy := "yrmn.yearmonth, yrmn.year, yrmn.mon"
	orderBy := "yrmn.yearmonth DESC"
	err := t.repo.FetchWithFullQuery(c, list, table, silect, joins, "", groupBy, orderBy, 0, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}

/* func (t *ExtController) Statement(c echo.Context) error {
	duration := c.Param("duration")
	date_condition := fmt.Sprintf(" event_date ", duration)
	list := &[]models.Account{}
	err := t.repo.FetchWithFullQuery(c, list, table, silect, joins, where, groupBy, orderBy, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, list)
} */
