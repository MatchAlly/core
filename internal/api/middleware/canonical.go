package middleware

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

func CanonicalLogger(log *zap.SugaredLogger) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		start := time.Now()

		next(ctx)

		duration := time.Since(start)

		log.Infow("request",
			"method", ctx.Method(),
			"path", ctx.URL().Path,
			"status", ctx.Status(),
			"duration_ms", duration.Milliseconds(),
			"ip", ctx.RemoteAddr(),
			"error", ctx.Context().Value("error"),
		)
	}
}
