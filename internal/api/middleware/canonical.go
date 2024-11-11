package middleware

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func CanonicalLogger(log *zap.SugaredLogger) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		start := time.Now()

		next(ctx)

		duration := time.Since(start)

		log.Infow("request",
			"requestID", middleware.GetReqID(ctx.Context()),
			"method", ctx.Method(),
			"path", ctx.URL().Path,
			"status", ctx.Status(),
			"duration_ms", duration.Milliseconds(),
			"ip", ctx.RemoteAddr(),
		)
	}
}
