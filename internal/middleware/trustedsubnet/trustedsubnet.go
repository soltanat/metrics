package trustedsubnet

import (
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
)

func Middleware(trustedSubnet string) (echo.MiddlewareFunc, error) {
	_, ipNet, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		return nil, err
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			realIPStr := c.RealIP()
			if realIPStr == "" {
				return c.NoContent(http.StatusForbidden)
			}

			realIP := net.ParseIP(realIPStr)
			if realIP == nil {
				return c.NoContent(http.StatusForbidden)
			}

			if ipNet.Contains(realIP) {
				return next(c)
			}

			return c.NoContent(http.StatusForbidden)
		}
	}, nil
}
