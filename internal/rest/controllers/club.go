package controllers

import (
	"core/internal/club"
	"core/internal/rest/handlers"
	"core/internal/rest/helpers"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type createClubRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *Handlers) CreateClub(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[createClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	userId, err := strconv.ParseUint(c.Claims.Subject, 10, 64)
	if err != nil {
		h.logger.Error("failed to parse userId",
			"error", err)
		return echo.ErrInternalServerError
	}

	_, err = h.clubService.CreateClub(ctx, req.Name, uint(userId))
	if err != nil {
		h.logger.Error("failed to create Club",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusCreated)
}

type deleteClubRequest struct {
	ClubId uint `json:"clubId" validate:"required,gt=0"`
}

func (h *Handlers) DeleteClub(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[deleteClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.clubService.DeleteClub(ctx, req.ClubId); err != nil {
		h.logger.Error("failed to delete Club",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

type updateClubRequest struct {
	ClubId uint   `json:"clubId" validate:"required,gt=0"`
	Name   string `json:"name" validate:"required"`
}

func (h *Handlers) UpdateClub(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[updateClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.clubService.UpdateClub(ctx, req.ClubId, req.Name); err != nil {
		h.logger.Error("failed to update Club",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

type updateUserRoleRequest struct {
	ClubId uint      `param:"clubId" validate:"required,gt=0"`
	UserId uint      `param:"userId" validate:"required,gt=0"`
	Role   club.Role `json:"role" validate:"required"`
}

func (h *Handlers) UpdateUserRole(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[updateUserRoleRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.clubService.UpdateUserRole(ctx, req.UserId, req.ClubId, req.Role); err != nil {
		h.logger.Error("failed to update user role",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

type getUsersInClubRequest struct {
	ClubId uint `query:"clubId" validate:"required,gt=0"`
}

type userInClub struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type getUsersInClubResponse struct {
	Users []userInClub `json:"users"`
}

func (h *Handlers) GetUsersInClub(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[getUsersInClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	userIds, err := h.clubService.GetUserIdsInClub(ctx, req.ClubId)
	if err != nil {
		h.logger.Error("failed to get userIds in Club",
			"error", err)
		return echo.ErrInternalServerError
	}

	users, err := h.userService.GetUsers(ctx, userIds)
	if err != nil {
		h.logger.Error("failed to get users",
			"error", err)
		return echo.ErrInternalServerError
	}

	respUsers := make([]userInClub, len(users))
	for i, u := range users {
		respUsers[i] = userInClub{
			Id:    u.Id,
			Name:  u.Name,
			Email: u.Email,
		}
	}

	resp := getUsersInClubResponse{
		Users: respUsers,
	}

	return c.JSON(http.StatusOK, resp)
}

type userInvitesClub struct {
	Id     uint   `json:"id"`
	ClubId uint   `json:"clubId"`
	Name   string `json:"name"`
}

type userInvitesResponse struct {
	Invites []userInvitesClub `json:"invites"`
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

	invites := make([]userInvitesClub, len(clubs))
	for i, c := range clubs {
		invites[i] = userInvitesClub{
			Id:     clubUsers[i].Id,
			ClubId: c.Id,
			Name:   c.Name,
		}
	}

	resp := userInvitesResponse{
		Invites: invites,
	}

	return c.JSON(http.StatusOK, resp)
}

type inviteUsersToClubRequest struct {
	ClubId uint     `json:"clubId"`
	Emails []string `json:"emails"`
}

func (h *Handlers) InviteUsersToClub(c handlers.AuthenticatedContext) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[inviteUsersToClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	users, err := h.userService.GetUsersByEmails(ctx, req.Emails)
	if err != nil {
		h.logger.Error("failed to get users by email",
			"error", err)
		return echo.ErrInternalServerError
	}

	userIds := make([]uint, len(users))
	for i, u := range users {
		userIds[i] = u.Id
	}

	if err := h.clubService.InviteToClub(ctx, userIds, req.ClubId); err != nil {
		h.logger.Error("failed to invite users to club",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
