package handlers

import (
	"context"
	"core/internal/member"

	"github.com/danielgtaylor/huma/v2"
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

func (h *Handler) GetMemberships(ctx context.Context, req *struct{}) (*getMembershipsResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	memberships, err := h.memberService.GetUserMemberships(ctx, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to get memberships, try again later")
	}

	clubIDs := make([]int, len(memberships))
	for i, m := range memberships {
		clubIDs[i] = m.ClubID
	}

	clubs, err := h.clubService.GetClubs(ctx, clubIDs)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to get clubs, try again later")
	}

	resp := &getMembershipsResponse{}
	mappedMemberships := make([]getMembershipsResponseClub, len(memberships))
	for i, m := range memberships {
		mappedMemberships[i] = getMembershipsResponseClub{
			ID:   m.ID,
			Name: clubs[i].Name,
		}
	}

	return resp, nil
}

type createClubRequest struct {
	Name string `json:"name" validate:"required"`
}

type createClubResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateClub(ctx context.Context, req *createClubRequest) (*createClubResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	clubID, err := h.clubService.CreateClub(ctx, req.Name, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to create club, try again later")
	}

	resp := &createClubResponse{
		ID:   clubID,
		Name: req.Name,
	}

	return resp, nil
}

type deleteClubRequest struct {
	ClubID int `json:"clubId" validate:"required,gt=0"`
}

func (h *Handler) DeleteClub(ctx context.Context, req *deleteClubRequest) (*struct{}, error) {
	if err := h.clubService.DeleteClub(ctx, req.ClubID); err != nil {
		return nil, huma.Error500InternalServerError("failed to delete club, try again later")
	}

	return nil, nil
}

type updateClubRequest struct {
	ClubID int    `json:"clubId" validate:"required,gt=0"`
	Name   string `json:"name" validate:"required"`
}

type updateClubResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) UpdateClub(ctx context.Context, req *updateClubRequest) (*updateClubResponse, error) {
	if err := h.clubService.UpdateClub(ctx, req.ClubID, req.Name); err != nil {
		return nil, huma.Error500InternalServerError("failed to update club, try again later")
	}

	resp := &updateClubResponse{
		ID:   req.ClubID,
		Name: req.Name,
	}

	return resp, nil
}

type updateMemberRoleRequest struct {
	MemberID int         `param:"clubId" validate:"required,gt=0"`
	Role     member.Role `json:"role" validate:"required"`
}

func (h *Handler) UpdateMemberRole(ctx context.Context, req *updateMemberRoleRequest) (*struct{}, error) {
	if err := h.memberService.UpdateRole(ctx, req.MemberID, req.Role); err != nil {
		return nil, huma.Error500InternalServerError("failed to update role, try again later")
	}

	return nil, nil
}

type getMembersInClubRequest struct {
	ClubId int `query:"clubId" validate:"required,gt=0"`
}

type getMembersInClubResponse struct {
	Members []membersInClub `json:"members"`
}

type membersInClub struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func (h *Handler) GetMembersInClub(ctx context.Context, req *getMembersInClubRequest) (*getMembersInClubResponse, error) {
	members, err := h.memberService.GetMembersInClub(ctx, req.ClubId)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to get members, try again later")
	}

	membersResponse := make([]membersInClub, len(members))
	for i, m := range members {
		membersResponse[i] = membersInClub{
			Id:   m.ID,
			Name: m.DisplayName,
			Role: string(m.Role),
		}
	}

	resp := &getMembersInClubResponse{
		Members: membersResponse,
	}

	return resp, nil
}

type removeUserFromClubRequest struct {
	MemberId int `param:"memberId" validate:"required,gt=0"`
}

func (h *Handler) RemoveMemberFromClub(ctx context.Context, req *removeUserFromClubRequest) (*struct{}, error) {
	if err := h.memberService.DeleteMembership(ctx, req.MemberId); err != nil {
		return nil, huma.Error500InternalServerError("failed to remove member from club, try again later")
	}

	return nil, nil
}
