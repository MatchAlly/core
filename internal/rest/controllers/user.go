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

func (h *Handlers) RemoveUserFromClub(c handlers.AuthenticatedContext) error {
	type request struct {
		UserId uint `param:"userId" validate:"required,gt=0"`
		ClubId uint `param:"clubId" validate:"required,gt=0"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[request](c)
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

func (h *Handlers) AddVirtualUserToClub(c handlers.AuthenticatedContext) error {
	type request struct {
		Name string `json:"name" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[request](c)
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

func (h *Handlers) RespondToInvite(c handlers.AuthenticatedContext) error {
	/*
		type request struct {
			InviteId uint `param:"inviteId" validate:"required,gt=0"`
			Accept   bool `json:"accept" validate:"required"`
		}

		ctx := c.Request().Context()

		req, err := helpers.Bind[request](c)
		if err != nil {
			return echo.ErrBadRequest
		}

		if err := h.clubService.RespondToInvite(ctx, req.InviteId, req.Accept); err != nil {
			h.logger.Error("failed to respond to invite",
				"error", err)
			return echo.ErrInternalServerError
		}
	*/
	return c.NoContent(http.StatusOK)
}
