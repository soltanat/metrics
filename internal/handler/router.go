package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes(h *Handlers) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	r := e.Router()
	r.Add(echo.GET, "/", h.GetList)
	r.Add(echo.GET, "/value/:metricType/:metricName", h.Get)
	r.Add(echo.POST, "/update/:metricType/:metricName/:metricValue", h.Store)

	return e
}
