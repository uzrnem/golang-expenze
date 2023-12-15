package loader

import (
	"expensez/src/controller"
	"fmt"

	pkgErr "github.com/uzrnem/go/errors"
	"github.com/uzrnem/go/rdb"
	"github.com/uzrnem/go/utils"
	"github.com/uzrnem/go/validator"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrorMap = map[string]error{
		"Name_required":  pkgErr.New(pkgErr.VALIDATION_ERROR, "Name is required"),
		"Email_required": pkgErr.New(pkgErr.VALIDATION_ERROR, "Email is required"),
		"Email_email":    pkgErr.New(pkgErr.VALIDATION_ERROR, "Email is invalid"),
	}
	MysqlDB *gorm.DB
)

type Loader struct {
	load func() error
	err  string
}

func Load() error {

	err := validator.Load(ErrorMap)
	if err != nil {
		return pkgErr.New(pkgErr.INTERNAL_ERROR, "Validator: failed to load")
	}

	err = DBLoad()
	if err != nil {
		return pkgErr.New(pkgErr.INTERNAL_ERROR, "MysqlDB: failed to load")
	}

	err = rdb.Load(MysqlDB)
	if err != nil {
		return pkgErr.New(pkgErr.INTERNAL_ERROR, "BaseRepository: failed to load")
	}

	loaders := []Loader{
		{load: controller.TagLoad, err: "TagController: failed to load"},
		{load: controller.TransactionTypeLoad, err: "TransactionTypeController: failed to load"},
		{load: controller.AccountTypeLoad, err: "AccountTypeController: failed to load"},
		{load: controller.AccountLoad, err: "AccountController: failed to load"},
		{load: controller.ActivityLoad, err: "ActivityController: failed to load"},
		{load: controller.PassbookLoad, err: "PassbookController: failed to load"},
		{load: controller.StatementLoad, err: "StatementController: failed to load"},
		{load: controller.ExtendedLoad, err: "ExtendedController: failed to load"},
		{load: controller.SubscriptionLoad, err: "SubscriptionController: failed to load"},
		{load: controller.StockLoad, err: "StockController: failed to load"},
	}

	for _, loader := range loaders {
		err = loader.load()
		if err != nil {
			return pkgErr.New(pkgErr.INTERNAL_ERROR, loader.err)
		}
	}
	return nil
}

func DBLoad() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.ReadEnvOrDefault("EXP_CONFIG_USER_NAME", "root"),
		utils.ReadEnvOrDefault("EXP_CONFIG_PASS_NAME", "root"),
		utils.ReadEnvOrDefault("EXP_CONFIG_HOST_NAME", "127.0.0.1"),
		utils.ReadEnvOrDefault("EXP_CONFIG_PORT_NUMB", "3306"),
		utils.ReadEnvOrDefault("EXP_CONFIG_DATABASE_NAME", "expense"),
	)
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	MysqlDB = db
	return nil
}
