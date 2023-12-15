package main

import (
	"expensez/src/loader"
	"expensez/src/routes"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/uzrnem/go/utils"
	"golang.org/x/sync/errgroup"
)

var gr errgroup.Group

func main() {
	var err error

	err = loader.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	e := echo.New()
	routes.SetupMiddlerware(e)
	routes.SetupController(e)

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
