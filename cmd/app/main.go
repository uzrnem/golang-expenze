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
	routes.SetupRoutes(e, "/tag", controller.TagCtrl)

	//utils.ReadEnvOrDefault("user", "root")
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
