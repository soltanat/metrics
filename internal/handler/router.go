package handler

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ziflex/lecho/v3"

	"github.com/soltanat/metrics/internal/logger"
)

func SetupRoutes(h *Handlers) *echo.Echo {
	l := logger.Get()

	e := echo.New()

	e.HideBanner = true

	e.Logger = lecho.New(l)

	e.Pre(middleware.AddTrailingSlash())

	e.Use(middleware.Decompress())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			ct := c.Request().Header.Get("Content-Type")
			if strings.Contains("application/json", ct) {
				return false
			}
			if strings.Contains("text/plain", ct) {
				return false
			}
			return true
		},
		Level:     -1,
		MinLength: 0,
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:          true,
		LogStatus:       true,
		LogLatency:      true,
		LogMethod:       true,
		LogResponseSize: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			l.Info().
				Str("URI", v.URI).
				Str("method", v.Method).
				Int("status", v.Status).
				Dur("latency", v.Latency).
				Int64("response_size", v.ResponseSize).
				Msg("request processed")
			return nil
		},
	}))
	e.Use(middleware.Recover())

	r := e.Router()
	r.Add(echo.GET, "/", h.GetList)
	r.Add(echo.GET, "/value/:metricType/:metricName/", h.Get)
	r.Add(echo.POST, "/update/:metricType/:metricName/:metricValue/", h.Store)
	r.Add(echo.POST, "/update/", h.StoreMetrics)
	r.Add(echo.POST, "/value/", h.Value)

	return e
}
