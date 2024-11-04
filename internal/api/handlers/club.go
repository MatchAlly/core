package handlers

import (
	"core/internal/api/helpers"
	"core/internal/member"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type getMembershipsResponse struct {
	Clubs    []getMembershipsResponseClub    `json:"clubs"`
	Invites  []getMembershipsResponseInvite  `json:"invites"`
	Requests []getMembershipsResponseRequest `json:"requests"`
}

type getMembershipsResponseClub struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type getMembershipsResponseInvite struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type getMembershipsResponseRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetMemberships(c helpers.AuthContext) error {
	ctx := c.Request().Context()

	memberships, err := h.memberService.GetUserMemberships(ctx, c.UserID)
	if err != nil {
		h.l.Error("failed to get memberships", zap.Error(err))
		return echo.ErrInternalServerError
	}

	clubIDs := make([]int, len(memberships))
	for i, m := range memberships {
		clubIDs[i] = m.ClubID
	}

	clubs, err := h.clubService.GetClubs(ctx, clubIDs)
	if err != nil {
		h.l.Error("failed to get clubs", zap.Error(err))
		return echo.ErrInternalServerError
	}

	response := getMembershipsResponse{}
	mappedMemberships := make([]getMembershipsResponseClub, len(memberships))
	for i, m := range memberships {
		mappedMemberships[i] = getMembershipsResponseClub{
			ID:   m.ID,
			Name: clubs[i].Name,
		}
	}

	return c.JSON(http.StatusOK, response)
}

type createClubRequest struct {
	Name string `json:"name" validate:"required"`
}

type createClubResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateClub(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[createClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	clubID, err := h.clubService.CreateClub(ctx, req.Name, c.UserID)
	if err != nil {
		h.l.Error("failed to create Club", zap.Error(err))
		return echo.ErrInternalServerError
	}

	response := createClubResponse{
		ID:   clubID,
		Name: req.Name,
	}

	return c.JSON(http.StatusCreated, response)
}

type deleteClubRequest struct {
	ClubID int `json:"clubId" validate:"required,gt=0"`
}

func (h *Handler) DeleteClub(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[deleteClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.clubService.DeleteClub(ctx, req.ClubID); err != nil {
		h.l.Error("failed to delete Club", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

type updateClubRequest struct {
	ClubID int    `json:"clubId" validate:"required,gt=0"`
	Name   string `json:"name" validate:"required"`
}

func (h *Handler) UpdateClub(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[updateClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.clubService.UpdateClub(ctx, req.ClubID, req.Name); err != nil {
		h.l.Error("failed to update Club", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

type updateMemberRoleRequest struct {
	MemberID int         `param:"clubId" validate:"required,gt=0"`
	Role     member.Role `json:"role" validate:"required"`
}

func (h *Handler) UpdateMemberRole(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[updateMemberRoleRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.memberService.UpdateRole(ctx, req.MemberID, req.Role); err != nil {
		h.l.Error("failed to update member role", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

type getMembersInClubRequest struct {
	ClubId int `query:"clubId" validate:"required,gt=0"`
}

type membersInClub struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func (h *Handler) GetMembersInClub(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[getMembersInClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	members, err := h.memberService.GetMembersInClub(ctx, req.ClubId)
	if err != nil {
		h.l.Error("failed to get members in club", zap.Error(err))
		return echo.ErrInternalServerError
	}

	response := make([]membersInClub, len(members))
	for i, m := range members {
		response[i] = membersInClub{
			Id:   m.ID,
			Name: m.DisplayName,
			Role: string(m.Role),
		}
	}

	return c.JSON(http.StatusOK, response)
}

type removeUserFromClubRequest struct {
	MemberId int `param:"memberId" validate:"required,gt=0"`
}

func (h *Handler) RemoveMemberFromClub(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[removeUserFromClubRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.memberService.DeleteMembership(ctx, req.MemberId); err != nil {
		h.l.Error("failed to delete membersihp", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
