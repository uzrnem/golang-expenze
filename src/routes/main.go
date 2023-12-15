package routes

import (
	"expensez/src/controller"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetupMiddlerware(e *echo.Echo) {
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
}

func SetupController(e *echo.Echo) {
	SetupRoutes(e, "/tags", controller.TagCtrl)
	SetupRoutes(e, "/transaction_types", controller.TransactionTypeCtrl)
	SetupRoutes(e, "/accounts", controller.AccountCtrl)
	SetupRoutes(e, "/activities", controller.ActivityCtrl)
	SetupRoutes(e, "/passbooks", controller.PassbookCtrl)
	SetupRoutes(e, "/statements", controller.StatementCtrl)
	SetupRoutes(e, "/account_types", controller.AccountTypeCtrl)
	SetupRoutes(e, "/subscriptions", controller.SubscriptionCtrl)
	SetupRoutes(e, "/stocks", controller.StockCtrl)
	e.GET("/accounts/type/:accountType", controller.ExtendedCtrl.FindAccountByType)
	e.GET("/tags/transactions/hits", controller.ExtendedCtrl.GetTagsByTranscationHits)
	e.GET("/accounts/chart/share", controller.ExtendedCtrl.BalanceSheet)
	e.GET("/statements/monthly/:duration", controller.ExtendedCtrl.Statement)
	e.GET("/expenses/:monyear", controller.ExtendedCtrl.ExpenseSheet)
}

func SetupRoutes(e *echo.Echo, path string, handlers controller.Controller) {
	e.POST(path, handlers.Create)
	e.DELETE(path+"/:id", handlers.Delete)
	e.GET(path+"/:id", handlers.Get)
	e.POST(path+"/:id", handlers.Update)
	e.GET(path, handlers.List)
}
