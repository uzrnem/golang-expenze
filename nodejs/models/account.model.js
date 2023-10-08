'user strict';
var config = require('./../db.config');

var Account = function(account) {
  this.name = account.name;
  this.account_type_id = account.account_type_id;
  this.amount = account.amount;
  this.is_frequent = account.is_frequent ? account.is_frequent : false;
  this.is_snapshot_disable = account.is_snapshot_disable ? account.is_snapshot_disable : false;
  this.is_closed = account.is_closed ? account.is_closed : false;
  this.created_at = new Date();
  this.updated_at = new Date();
};

Account.create = function(newAccount, result) {
  config.con.query("INSERT INTO accounts set ?", newAccount, function(err, res) {
    if (err) {
      console.error("error: ", err);
      result({error: err.sqlMessage}, null);
    } else {
      result(null, res.insertId);
    }
  });
};

Account.findById = function(id, result) {
  config.con.query("Select * from accounts where id = ? ", id, function(err, res) {
    if (err) {
      console.error("error: ", err);
      result({error: err.sqlMessage}, null);
    } else {
      result(null, res);
    }
  });
};

Account.findAll = function(allAccounts, result) {
  var condition = "";
  if ( allAccounts ) {
    condition = " ORDER BY is_closed, account_type_id, amount = 0 ASC, name ASC";
  } else {
    condition = " WHERE is_closed = 0 ORDER BY name ASC";
  }
  config.con.query("Select * from accounts " + condition, function(err, res) {
    if (err) {
      console.error("error: ", err);
      result({error: err.sqlMessage}, null);
    } else {
      result(null, res);
    }
  });
};

Account.frequent = function(result) {
  config.con.query("Select * from accounts WHERE is_closed = 0 AND is_frequent = 1 ORDER BY name ASC", function(err, res) {
    if (err) {
      console.error("error: ", err);
      result({error: err.sqlMessage}, null);
    } else {
      result(null, res);
    }
  });
};

Account.findByType = function(account_type, result) {
  var sql = "Select * from accounts WHERE account_type_id in ( " +
  "SELECT id FROM account_types where name = '"+ account_type +"')";
  config.con.query(sql, function(err, res) {
    if (err) {
      console.error("error: ", err);
      result({error: err.sqlMessage}, null);
    } else {
      result(null, res);
    }
  });
};

Account.update = function(id, account, result) {
  config.con.query("UPDATE accounts SET name=?,account_type_id=?,amount=?,is_frequent=?,is_snapshot_disable=?,is_closed=? WHERE id = ?",
    [account.name, account.account_type_id, account.amount
      , account.is_frequent, account.is_snapshot_disable, account.is_closed, id],
    function(err, res) {
      if (err) {
        console.error("error: ", err);
        result({error: err.sqlMessage}, null);
      } else {
        result(null, res);
      }
    });
};

Account.delete = function(id, result) {
  config.con.query("DELETE FROM accounts WHERE id = ?", [id], function(err, res) {
    if (err) {
      console.error("error: ", err);
      result({error: err.sqlMessage}, null);
    } else {
      result(null, res);
    }
  });
};

Account.share = function(result) {
  (async () => {
    try {
      const holding_balance = await config.query ("select " +
      "  t.name as 'Account', SUM(a.amount) as 'Amount per Account' " +
      " from accounts a " +
      " left join account_types t on a.account_type_id = t.id " +
      " where a.amount !=0 and a.is_snapshot_disable = 0 and a.is_closed != 1 " +
      " group by a.account_type_id order by t.name='Saving' desc, t.name='Credit' desc, t.name='Wallet' desc, " +
      " t.name='Stocks Equity' desc, t.name='Loan' desc, t.name='Mutual Funds' desc, t.name='Deposit' desc;");

      const account_balance = await config.query(
        " select a.name as account, t.name as type, a.amount as balance " +
        " from accounts a " +
        " left join account_types t on a.account_type_id = t.id " +
        " where a.amount !=0 and a.is_snapshot_disable = 0 and a.is_closed != 1 " +
        " order by t.name='Saving' desc, t.name='Credit' desc, t.name='Wallet' desc, " +
        " t.name='Deposit' desc, t.name='Loan' desc, t.name='Stocks Equity', a.name;");

      var ccBills = {'Balance' : 0}
      var total = 0.0
      var ccBill = 0
      var loan = 0.0
      console.debug(holding_balance)
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
      delete ccBills['CC Bill']
      console.debug('total: ', total, 'Loan: ', loan, 'cc bill: ', ccBill)
      cc_array = [['Account', 'Amount per Account']]
      cc_array.push(['Loan', loan])
      cc_array.push(['CC Bill', ccBill])
      cc_array.push(['Balance', total])
      //creditLimit = 600000-loan-ccBill-total
      //cc_array.push(['Credit Limit', creditLimit > 0 ? creditLimit: 0])
      result(null, { holding: holding_array, balance: account_balance, totalBalance: cc_array });
    } finally {
    }
  })()
};

Account.expenses = function(month, result) {
  (async () => {
    try {
      var condition = " AND EXTRACT(YEAR_MONTH FROM event_date) = " + month
      if (month == 0) {
        condition = ""
      }
      const holding_balance = await config.query (
        "SELECT COALESCE(sub.name, tag.name) as tag, SUM(act.amount) as amount " +
        " FROM `activities` as act " +
        " LEFT JOIN `tags` tag ON `tag`.`id` = `act`.`tag_id` " +
        " LEFT JOIN `tags` sub ON `sub`.`id` = `act`.`sub_tag_id` " +
        " WHERE `act`.`transaction_type_id` = 2 " + condition +
        " GROUP BY tag.name, sub.name ORDER BY SUM(act.amount) ASC"
      );

      const months = await config.query(
        "SELECT EXTRACT(YEAR_MONTH FROM event_date) as mon, MONTHNAME(event_date) as month, year(event_date) as year " +
        "FROM activities WHERE transaction_type_id = 2 " +
        "GROUP BY EXTRACT(YEAR_MONTH FROM event_date), MONTHNAME(event_date), year(event_date) " +
        "ORDER BY EXTRACT(YEAR_MONTH FROM event_date) DESC;");

      const expenseByAccount = await config.query(
        " SELECT a.name as account, SUM(act.amount) as expense FROM `activities` act " +
        " LEFT JOIN accounts a ON act.from_account_id = a.id " +
        " WHERE act.transaction_type_id = 2 " + condition +
        " GROUP BY a.id, a.name, act.from_account_id ORDER BY SUM(act.amount) DESC;"
      )
      var total = 0.0
      holding_array = [['Tag', 'Amount']]
      holding_balance.forEach((item, i) => {
        holding_array.push([item['tag'], item['amount']])
        total = total + item['amount'];
      });
      result(null, { holding: holding_array, expenses: total, months: months, expenseByAccount });
      return;
    } finally {
    }
  })()
};

module.exports = Account;
