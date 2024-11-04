package handlers

import (
	"core/internal/api/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type updateUserRequest struct {
	Email string `json:"userId" validate:"required,email"`
	Name  string `json:"clubId" validate:"required,min=1,max=255"`
}

func (h *Handler) UpdateUser(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[updateUserRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.userService.UpdateUser(ctx, c.UserID, req.Email, req.Name); err != nil {
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) DeleteUser(c helpers.AuthContext) error {
	ctx := c.Request().Context()

	if err := h.userService.DeleteUser(ctx, c.UserID); err != nil {
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
