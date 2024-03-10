package handlers

import (
	"core/internal/api/helpers"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type updateUserRequest struct {
	Email string `json:"userId" validate:"required,email"`
	Name  string `json:"clubId" validate:"required,min=1,max=255"`
}

func (h *Handler) UpdateUser(c helpers.AuthenticatedContext) error {
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

func (h *Handler) DeleteUser(c helpers.AuthenticatedContext) error {
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

type responseInvite struct {
	Id     uint   `json:"id"`
	ClubId uint   `json:"clubId"`
	Name   string `json:"name"`
}

type getUserInvitesResponse struct {
	Invites []responseInvite `json:"invites"`
}

func (h *Handler) GetUserInvites(c helpers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	userId, err := strconv.ParseUint(c.Claims.Subject, 10, 64)
	if err != nil {
		h.logger.Error("failed to parse userId",
			"error", err)
		return echo.ErrInternalServerError
	}

	invites, err := h.clubService.GetInvitesByUserId(ctx, uint(userId))
	if err != nil {
		h.logger.Error("failed to get user invites",
			"error", err)
		return echo.ErrInternalServerError
	}
	// hello
	clubIds := make([]uint, len(invites))
	for i, invite := range invites {
		clubIds[i] = invite.ClubId
	}

	clubs, err := h.clubService.GetClubs(ctx, clubIds)
	if err != nil {
		h.logger.Error("failed to get Clubs",
			"error", err)
		return echo.ErrInternalServerError
	}

	responseInvites := make([]responseInvite, len(clubs))
	for i, c := range clubs {
		responseInvites[i] = responseInvite{
			Id:     invites[i].Id,
			ClubId: c.Id,
			Name:   c.Name,
		}
	}

	resp := getUserInvitesResponse{
		Invites: responseInvites,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) RespondToInvite(c helpers.AuthenticatedContext) error {
	return nil // TODO
}
