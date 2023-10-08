package mysql

import (
	"expensez/pkg/utils"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	MysqlDB *gorm.DB
)

func Load() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.ReadEnvOrDefault("EXP_CONFIG_USER_NAME", "root"),
		utils.ReadEnvOrDefault("EXP_CONFIG_PASS_NAME", "root"),
		utils.ReadEnvOrDefault("EXP_CONFIG_HOST_NAME", "127.0.0.1"),
		utils.ReadEnvOrDefault("EXP_CONFIG_PORT_NUMB", "3306"),
		utils.ReadEnvOrDefault("EXP_CONFIG_DATABASE_NAME", "expense_dev"),
	)
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	MysqlDB = db
	return nil
}
