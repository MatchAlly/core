package middleware

import (
	"log/slog"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

func CanonicalLogger(log *slog.Logger) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		start := time.Now()

		next(ctx)

		duration := time.Since(start)

		log.Info("Request",
			"method", ctx.Method(),
			"path", ctx.URL().Path,
			"status", ctx.Status(),
			"duration_ms", duration.Milliseconds(),
			"ip", ctx.RemoteAddr(),
		)
	}
}
