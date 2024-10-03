package logger

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewBasicRequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:       true,
		LogError:        true,
		LogLatency:      true,
		LogURIPath:      true,
		LogURI:          true,
		LogResponseSize: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			x := fmt.Sprintf("%s %s", c.Request().Method, v.URI)

			if v.Status > 199 && v.Status < 300 {
				log.Info(x, "code", v.Status, "millis", v.Latency.Milliseconds(), "bytes", v.ResponseSize)
			} else {
				log.Error(x, "code", v.Status, "millis", v.Latency.Milliseconds(), "bytes", v.ResponseSize)
			}

			return nil
		},
	})
}
