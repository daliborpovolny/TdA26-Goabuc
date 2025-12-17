package middlewares

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

var LoggerMiddleware = middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	LogStatus:   true,
	LogURI:      true,
	LogError:    true,
	HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
	LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
		if v.Error == nil {
			Logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
			)
		} else {
			Logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("err", v.Error.Error()),
			)
		}
		return nil
	},
})
