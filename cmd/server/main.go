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
	e.GET("/value/:metricType/:metricName", h.GetGauge)
	e.GET("/value/:metricType/:metricName", h.GetCounter)

	e.POST("/update/:metricType/:metricName/:metricValue", h.StoreGauge)
	e.POST("/update/:metricType/:metricName/:metricValue", h.StoreCounter)

	e.Logger.Fatal(e.Start(":8080"))
}
