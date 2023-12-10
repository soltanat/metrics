package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/soltanat/metrics/internal/handler"
	"github.com/soltanat/metrics/internal/storage"
)

func main() {
	s := storage.NewMemStorage()
	h := handler.New(s)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", h.GetList)
	e.GET("/value/gauge/:metricName", h.GetGauge)
	e.GET("/value/counter/:metricName", h.GetCounter)

	e.POST("/update/gauge/:metricName/:metricValue", h.StoreGauge)
	e.POST("/update/counter/:metricName/:metricValue", h.StoreCounter)

	e.Logger.Fatal(e.Start(":8080"))
}
