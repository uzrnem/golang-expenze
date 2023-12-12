package main

import (
	"expensez/di"
	"expensez/pkg/utils"
	"expensez/src/controller"
	"expensez/src/routes"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"golang.org/x/sync/errgroup"
)

var gr errgroup.Group

func main() {
	var err error

	err = di.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	e := echo.New()
	routes.SetupMiddlerware(e)
	routes.SetupRoutes(e, "/tags", controller.TagCtrl)
	routes.SetupRoutes(e, "/transaction_types", controller.TransactionTypeCtrl)
	routes.SetupRoutes(e, "/accounts", controller.AccountCtrl)
	routes.SetupRoutes(e, "/activities", controller.ActivityCtrl)
	routes.SetupRoutes(e, "/passbooks", controller.PassbookCtrl)
	routes.SetupRoutes(e, "/statements", controller.StatementCtrl)
	routes.SetupRoutes(e, "/account_types", controller.AccountTypeCtrl)
	routes.SetupRoutes(e, "/subscriptions", controller.SubscriptionCtrl)
	routes.SetupRoutes(e, "/stocks", controller.StockCtrl)
	e.GET("/accounts/type/:accountType", controller.ExtendedCtrl.FindAccountByType)
	e.GET("/tags/transactions/hits", controller.ExtendedCtrl.GetTagsByTranscationHits)
	e.GET("/accounts/chart/share", controller.ExtendedCtrl.BalanceSheet)
	e.GET("/statements/monthly/:duration", controller.ExtendedCtrl.Statement)
	e.GET("/expenses/:monyear", controller.ExtendedCtrl.ExpenseSheet)

	e.Static("/", "public")
	port := utils.ReadEnvOrDefault("SERVER_PORT", "9001")
	mainServer := &http.Server{
		Addr:    ":" + port,
		Handler: e,
	}

	startServer(mainServer, port)
}

func startServer(mainServer *http.Server, port string) {
	gr.Go(func() error {
		return mainServer.ListenAndServe()
	})
	log.Println("Service Running at: " + port)
	if err := gr.Wait(); err != nil {
		log.Fatal(err)
	}
}
