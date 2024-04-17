package helpers

import (
	"context"

	"github.com/creasty/defaults"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func Bind[V any](ctx echo.Context) (*V, context.Context, error) {
	var data V
	if err := ctx.Bind(&data); err != nil {
		return nil, nil, errors.Wrap(err, "failed to bind request")
	}

	if err := ctx.Validate(&data); err != nil {
		return nil, nil, errors.Wrap(err, "request failed validation")
	}

	if err := defaults.Set(&data); err != nil {
		return nil, nil, errors.Wrap(err, "failed to set defaults")
	}

	return &data, ctx.Request().Context(), nil
}
