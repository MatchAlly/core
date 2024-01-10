package controllers

import (
	"core/internal/rest/handlers"
	"core/internal/rest/helpers"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) DeleteUser(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	userId, err := strconv.ParseUint(c.Claims.Subject, 10, 64)
	if err != nil {
		h.logger.Error("failed to parse userId",
			"error", err)
		return echo.ErrInternalServerError
	}

	if err := h.userService.DeleteUser(ctx, uint(userId)); err != nil {
		h.logger.Error("failed to delete user",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

type removeUserFromClubRequest struct {
	UserId uint `param:"userId" validate:"required,gt=0"`
	ClubId uint `param:"clubId" validate:"required,gt=0"`
}

func (h *Handlers) RemoveUserFromClub(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[removeUserFromClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.clubService.RemoveUserFromClub(ctx, req.UserId, req.ClubId); err != nil {
		h.logger.Error("failed to remove user from Club",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

type addVirtualUserToClubRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *Handlers) AddVirtualUserToClub(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[addVirtualUserToClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}
	if err := h.userService.CreateVirtualUser(ctx, req.Name); err != nil {
		h.logger.Error("failed to create virtual user",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
