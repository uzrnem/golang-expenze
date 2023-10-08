package controller

import ( "fmt"
	"strings"
	v "expensez/pkg/validator"
	"expensez/src/models"
	"expensez/src/repository"
	"net/http"
	"expensez/pkg/utils"
	//"strconv"

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
	where := "account_type_id in (SELECT id FROM account_types where name = '"+ account_type +"')"
	order := "name ASC"
	err := t.repo.FetchWithQuery(c, list, where, order)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}

type FullExtended12 struct {
	FromAccount string   `json:"from_account"`
	ToAccount string   `json:"to_account"`
	Tag string   `json:"tag"`
	SubTag string   `json:"sub_tag"`
	TransactionaType string   `json:"transaction_type"`
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
		"tags": nil,
		"sub_tags": nil,
		"tag_id": utils.StringToInt(tagID),
		"sub_tag_id": 0,
		"remarks": "",
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
		conditions = append(conditions, "tag_id = " + tagID)
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
		err = t.repo.FetchWithQuery(c, tagList, "tag_id = " + tagID, "")
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


/*
Account.share = function(result) {
	(async () => {
	  try {
		const holding_balance = await config.query ("select 
		t.name as name, SUM(a.amount) as amount 
		from accounts a 
		left join account_types t on a.account_type_id = t.id 
		where a.amount !=0 and a.is_snapshot_disable = 0 and a.is_closed != 1 
		group by a.account_type_id order by t.name='Saving' desc, t.name='Credit' desc, t.name='Wallet' desc, 
		t.name='Stocks Equity' desc, t.name='Loan' desc, t.name='Mutual Funds' desc, t.name='Deposit' desc;");
		holding_array = [['Account', 'Amount per Account']]
		holding_balance.forEach((item, i) => {
		  if (item['Account'] == 'Credit') {
			ccBill = 0 - item['Amount per Account']
			//holding_array.push([item['Account'], ccBill])
		  } else {
			holding_array.push([item['Account'], item['Amount per Account']])
		  }
		  if (item['Account'] == 'Deposit' || item['Account'] == 'Stocks Equity' || item['Account'] == 'Mutual Funds') {
  
		  } else if (item['Account'] == 'Loan') {
			loan = item['Amount per Account']
		  } else {
			total = total + item['Amount per Account']
		  }
		});


		
  
		const account_balance = await config.query(
		" select a.name as account, t.name as type, a.amount as amount 
		from accounts a 
		left join account_types t on a.account_type_id = t.id 
		where a.amount !=0 and a.is_snapshot_disable = 0 and a.is_closed != 1 
		order by t.name='Saving' desc, t.name='Credit' desc, t.name='Wallet' desc, 
		t.name='Deposit' desc, t.name='Loan' desc, t.name='Stocks Equity', a.name;");
  
		var ccBills = {'Balance' : 0}
		var total = 0.0
		var ccBill = 0
		var loan = 0.0
		console.debug(holding_balance)
		delete ccBills['CC Bill']
		console.debug('total: ', total, 'Loan: ', loan, 'cc bill: ', ccBill)
		cc_array = [['Account', 'Amount per Account']]
		cc_array.push(['Loan', loan])
		cc_array.push(['CC Bill', ccBill])
		cc_array.push(['Balance', total])
		result(null, { holding: holding_array, balance: account_balance, totalBalance: cc_array });
	  } finally {
	  }
	})()
  };*/