package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func CanonicalLogger(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			res := c.Response()

			err := next(c)

			latency := float64(time.Since(start).Milliseconds())

			logFunc := logger.Infow
			msg := "Request completed"

			switch {
			case res.Status >= 500:
				logFunc = logger.Errorw
				msg = "Server error"
			case res.Status >= 400:
				logFunc = logger.Warnw
				msg = "Client error"
			case res.Status >= 300:
				msg = "Redirect"
			}

			logFunc(msg,
				"method", req.Method,
				"path", req.URL.Path,
				"status", res.Status,
				"latency_ms", latency,
				"content_length", res.Size,
				"user_agent", req.UserAgent(),
				"remote_ip", c.RealIP(),
				"request_id", req.Header.Get(echo.HeaderXRequestID),
				"query", req.URL.RawQuery,
			)

			if err != nil {
				logger.Errorw("Request error",
					"error", err.Error(),
					"path", req.URL.Path,
					"method", req.Method,
				)
			}

			return err
		}
	}
}
