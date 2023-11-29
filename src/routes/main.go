package routes

import (
	"fmt"
	"expensez/src/controller"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetupMiddlerware(e *echo.Echo) {
	e.Pre(middleware.RemoveTrailingSlash())
	//e.Use(Counter())
	//e.Use(middleware.Logger())
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
	e.POST(path+"/:id", handlers.Update)
	e.GET(path, handlers.List)
}

func Counter() echo.MiddlewareFunc {
	logMap := map[string]map[string]int{}
	return func (next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			next(c)
			req := c.Request()
			url := req.URL
			path := url.Path
	
			if logMap[req.Method] == nil {
				logMap[req.Method] = map[string]int{}
			}
			logMap[req.Method][path] = logMap[req.Method][path] + 1
			fmt.Printf("Counter: %+v\n", logMap)
			return nil
		}
	}
}