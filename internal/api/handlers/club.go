package handlers

import (
	"core/internal/api/helpers"
	"core/internal/club"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type createClubRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *Handler) CreateClub(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[createClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	userId, err := strconv.ParseUint(c.Claims.Subject, 10, 64)
	if err != nil {
		h.l.Error("failed to parse userId", zap.Error(err))
		return echo.ErrInternalServerError
	}

	_, err = h.clubService.CreateClub(ctx, req.Name, uint(userId))
	if err != nil {
		h.l.Error("failed to create Club", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusCreated)
}

type deleteClubRequest struct {
	ClubId uint `json:"clubId" validate:"required,gt=0"`
}

func (h *Handler) DeleteClub(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[deleteClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.clubService.DeleteClub(ctx, req.ClubId); err != nil {
		h.l.Error("failed to delete Club", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

type updateClubRequest struct {
	ClubId uint   `json:"clubId" validate:"required,gt=0"`
	Name   string `json:"name" validate:"required"`
}

func (h *Handler) UpdateClub(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[updateClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.clubService.UpdateClub(ctx, req.ClubId, req.Name); err != nil {
		h.l.Error("failed to update Club", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

type updateMemberRoleRequest struct {
	MemberId uint      `param:"clubId" validate:"required,gt=0"`
	Role     club.Role `json:"role" validate:"required"`
}

func (h *Handler) UpdateMemberRole(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[updateMemberRoleRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.clubService.UpdateMemberRole(ctx, req.MemberId, req.Role); err != nil {
		h.l.Error("failed to update user role", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

type getMembersInClubRequest struct {
	ClubId uint `query:"clubId" validate:"required,gt=0"`
}

type membersInClub struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *Handler) GetMembersInClub(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[getMembersInClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	members, err := h.clubService.GetMembers(ctx, req.ClubId)
	if err != nil {
		h.l.Error("failed to get members in club", zap.Error(err))
		return echo.ErrInternalServerError
	}

	userIds := make([]uint, len(members))
	for i, m := range members {
		userIds[i] = m.UserId
	}

	users, err := h.userService.GetUsers(ctx, userIds)
	if err != nil {
		h.l.Error("failed to get users", zap.Error(err))
		return echo.ErrInternalServerError
	}

	resp := make([]membersInClub, len(users))
	for i, u := range users {
		resp[i] = membersInClub{
			Id:    members[i].Model.ID,
			Name:  u.Name,
			Email: u.Email,
			Role:  string(members[i].Role),
		}
	}

	return c.JSON(http.StatusOK, resp)
}

type removeUserFromClubRequest struct {
	MemberId uint `param:"memberId" validate:"required,gt=0"`
}

func (h *Handler) RemoveUserFromClub(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[removeUserFromClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.clubService.DeleteMember(ctx, req.MemberId); err != nil {
		h.l.Error("failed to delete member", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
