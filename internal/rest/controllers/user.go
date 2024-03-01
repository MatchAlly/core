package controllers

import (
	"core/internal/rest/handlers"
	"core/internal/rest/helpers"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type updateUserRequest struct {
	Email string `json:"userId" validate:"required,email"`
	Name  string `json:"clubId" validate:"required,min=1,max=255"`
}

func (h *Handlers) UpdateUser(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[updateUserRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	userId, err := strconv.ParseUint(c.Claims.Subject, 10, 64)
	if err != nil {
		h.logger.Error("failed to parse userId",
			"error", err)
		return echo.ErrInternalServerError
	}

	if err := h.userService.UpdateUser(ctx, uint(userId), req.Email, req.Name); err != nil {
		h.logger.Error("failed to delete user",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

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

type invite struct {
	Id     uint   `json:"id"`
	ClubId uint   `json:"clubId"`
	Name   string `json:"name"`
}

type getUserInvitesResponse struct {
	Invites []invite `json:"invites"`
}

func (h *Handlers) GetUserInvites(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	userId, err := strconv.ParseUint(c.Claims.Subject, 10, 64)
	if err != nil {
		h.logger.Error("failed to parse userId",
			"error", err)
		return echo.ErrInternalServerError
	}

	clubUsers, err := h.clubService.GetInvitesByUserId(ctx, uint(userId))
	if err != nil {
		h.logger.Error("failed to get user invites",
			"error", err)
		return echo.ErrInternalServerError
	}

	clubIds := make([]uint, len(clubUsers))
	for i, clubUser := range clubUsers {
		clubIds[i] = clubUser.ClubId
	}

	clubs, err := h.clubService.GetClubs(ctx, clubIds)
	if err != nil {
		h.logger.Error("failed to get Clubs",
			"error", err)
		return echo.ErrInternalServerError
	}

	invites := make([]invite, len(clubs))
	for i, c := range clubs {
		invites[i] = invite{
			Id:     clubUsers[i].Id,
			ClubId: c.Id,
			Name:   c.Name,
		}
	}

	resp := getUserInvitesResponse{
		Invites: invites,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handlers) RespondToInvite(c handlers.AuthenticatedContext) error {
	return nil // TODO
}
