package controller

import (
	"expensez/src/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	repository "github.com/uzrnem/go/rdb"
	"github.com/uzrnem/go/utils"
)

var (
	ActivityCtrl Controller
)

type ActivityController struct {
	BaseController
}

func ActivityLoad() error {
	ActivityCtrl = &ActivityController{
		BaseController: BaseController{
			GetModel: func(id string) (any, error) {
				modal := &models.Activity{}
				if id != "" {
					ID, err := strconv.Atoi(id)
					if err != nil {
						return nil, err
					}
					modal.ID = uint(ID)
				}
				return modal, nil
			},
			GetList: func() any {
				return &[]models.Activity{}
			},
		},
	}
	return nil
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

	list := &[]models.FullActivity{}
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
	err := repository.Repo.Builder().Table(table).Select(silect).Join(joins).Where(where).Order(orderBy).Limit(limit).Offset(offset).Exec(list)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var count, sum float64
	err = repository.Repo.Builder().Table(table).Select("COUNT(act.id) as count").Where(utils.TrimSpace(where)).Exec(&count)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if count < 100 {
		err = repository.Repo.Builder().Table(table).Select("SUM(act.amount) as sum").Where(utils.TrimSpace(where)).Exec(&sum)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	return c.JSON(http.StatusOK, map[string]any{
		"list":  list,
		"total": count,
		"sum":   sum,
	})
}
