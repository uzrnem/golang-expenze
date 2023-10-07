package routes

import (
	"expensez/src/controller"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetupMiddlerware(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
}

func SetupRoutes(e *echo.Echo, path string, handlers controller.Controller) {
	e.POST(path, handlers.Create)
	e.DELETE(path+"/:id", handlers.Delete)
	e.GET(path+"/:id", handlers.Get)
	e.PUT(path+"/:id", handlers.Update)
	e.GET(path, handlers.List)
}
