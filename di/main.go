package di

import (
	"expensez/config/mysql"
	pkgErr "expensez/pkg/errors"
	"expensez/pkg/validator"
	"expensez/src/controller"
	"expensez/src/repository"
)

var (
	ErrorMap = map[string]error{
		"Name_required":  pkgErr.New(pkgErr.VALIDATION_ERROR, "Name is required"),
		"Email_required": pkgErr.New(pkgErr.VALIDATION_ERROR, "Email is required"),
		"Email_email":    pkgErr.New(pkgErr.VALIDATION_ERROR, "Email is invalid"),
	}
)

type Loader struct {
	load func() error
	err  string
}

func Load() error {

	loaders := []Loader{
		{load: mysql.Load, err: "MysqlDB: failed to load"},
		{load: repository.Load, err: "BaseRepository: failed to load"},
		{load: controller.Load, err: "TagController: failed to load"},
	}

	err := validator.Load(ErrorMap)
	if err != nil {
		return pkgErr.New(pkgErr.INTERNAL_ERROR, "Validator: failed to load")
	}

	for _, loader := range loaders {
		err = loader.load()
		if err != nil {
			return pkgErr.New(pkgErr.INTERNAL_ERROR, loader.err)
		}
	}
	return nil
}
