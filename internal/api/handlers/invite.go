package handlers

import (
	"core/internal/api/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type inviteUsersToClubRequest struct {
	Emails []string `json:"emails"`
	ClubId uint     `json:"clubId"`
}

func (h *Handler) InviteUsersToClub(c helpers.AuthContext) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[inviteUsersToClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	users, err := h.userService.GetUsersByEmails(ctx, req.Emails)
	if err != nil {
		h.l.Error("failed to get users by email", zap.Error(err))
		return echo.ErrInternalServerError
	}

	userIds := make([]uint, len(users))
	for i, u := range users {
		userIds[i] = u.Model.ID
	}

	if err := h.inviteService.CreateInvites(ctx, userIds, req.ClubId); err != nil {
		h.l.Error("failed to invite users to club", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
