package controllers

import (
	"core/internal/club"
	"core/internal/rest/handlers"
	"core/internal/rest/helpers"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) CreateClub(c handlers.AuthenticatedContext) error {
	type request struct {
		Name string `json:"name" validate:"required"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[request](c)
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

func (h *Handlers) DeleteClub(c handlers.AuthenticatedContext) error {
	type request struct {
		ClubId uint `json:"clubId" validate:"required,gt=0"`
	}

	req, err := helpers.Bind[request](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	ctx := c.Request().Context()

	if err := h.clubService.DeleteClub(ctx, req.ClubId); err != nil {
		h.logger.Error("failed to delete Club",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) UpdateClub(c handlers.AuthenticatedContext) error {
	type request struct {
		ClubId uint   `json:"clubId" validate:"required,gt=0"`
		Name   string `json:"name" validate:"required"`
	}

	req, err := helpers.Bind[request](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	ctx := c.Request().Context()

	if err := h.clubService.UpdateClub(ctx, req.ClubId, req.Name); err != nil {
		h.logger.Error("failed to update Club",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) UpdateUserRole(c handlers.AuthenticatedContext) error {
	type request struct {
		ClubId uint      `param:"clubId" validate:"required,gt=0"`
		UserId uint      `param:"userId" validate:"required,gt=0"`
		Role   club.Role `json:"role" validate:"required"`
	}

	req, err := helpers.Bind[request](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	ctx := c.Request().Context()

	if err := h.clubService.UpdateUserRole(ctx, req.UserId, req.ClubId, req.Role); err != nil {
		h.logger.Error("failed to update user role",
			"error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) GetUsersInClub(c handlers.AuthenticatedContext) error {
	type request struct {
		ClubId uint `query:"clubId" validate:"required,gt=0"`
	}

	type responseUser struct {
		Id    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	type response struct {
		Users []responseUser `json:"users"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[request](c)
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

	respUsers := make([]responseUser, len(users))
	for i, u := range users {
		respUsers[i] = responseUser{
			Id:    u.Id,
			Name:  u.Name,
			Email: u.Email,
		}
	}

	resp := response{
		Users: respUsers,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handlers) GetUserInvites(c handlers.AuthenticatedContext) error {
	type responseClub struct {
		Id     uint   `json:"id"`
		ClubId uint   `json:"club_id"`
		Name   string `json:"name"`
	}

	type response struct {
		Invites []responseClub `json:"invites"`
	}

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

	invites := make([]responseClub, len(clubs))
	for i, c := range clubs {
		invites[i] = responseClub{
			Id:     clubUsers[i].Id,
			ClubId: c.Id,
			Name:   c.Name,
		}
	}

	resp := response{
		Invites: invites,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handlers) InviteUsersToClub(c handlers.AuthenticatedContext) error {
	type request struct {
		ClubId uint     `json:"club_id"`
		Emails []string `json:"emails"`
	}

	ctx := c.Request().Context()

	req, err := helpers.Bind[request](c)
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
