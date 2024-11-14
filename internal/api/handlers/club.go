package handlers

import (
	"context"
	"core/internal/member"

	"github.com/danielgtaylor/huma/v2"
)

type getMembershipsRequest struct {
	UserID int `path:"userId" minimum:"1"`
}

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

func (h *Handler) GetMemberships(ctx context.Context, req *getMembershipsRequest) (*getMembershipsResponse, error) {
	memberships, err := h.memberService.GetUserMemberships(ctx, req.UserID)
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
	Name string `json:"name" minLength:"2" maxLength:"64"`
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
	ClubID int `json:"clubId"  minimum:"1"`
}

func (h *Handler) DeleteClub(ctx context.Context, req *deleteClubRequest) (*struct{}, error) {
	if err := h.clubService.DeleteClub(ctx, req.ClubID); err != nil {
		return nil, huma.Error500InternalServerError("failed to delete club, try again later")
	}

	return nil, nil
}

type updateClubRequest struct {
	ClubID int    `json:"clubId" minimum:"1"`
	Name   string `json:"name" minLength:"2" maxLength:"64"`
}

type updateClubResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) UpdateClub(ctx context.Context, req *updateClubRequest) (*updateClubResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	ok, err := h.authZService.IsAdmin(ctx, userID, req.ClubID)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to update this club")
	}

	if err := h.clubService.UpdateClub(ctx, req.ClubID, req.Name); err != nil {
		return nil, huma.Error500InternalServerError("failed to update club")
	}

	resp := &updateClubResponse{
		ID:   req.ClubID,
		Name: req.Name,
	}

	return resp, nil
}

type updateMemberRoleRequest struct {
	MemberID int         `param:"clubId" minimum:"1"`
	Role     member.Role `json:"role" enum:"ADMIN,MANAGER,MEMBER"`
}

func (h *Handler) UpdateMemberRole(ctx context.Context, req *updateMemberRoleRequest) (*struct{}, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	ok, err := h.authZService.IsAdmin(ctx, userID, req.MemberID)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to update member role")
	}

	if err := h.memberService.UpdateRole(ctx, req.MemberID, req.Role); err != nil {
		return nil, huma.Error500InternalServerError("failed to update role, try again later")
	}

	return nil, nil
}

type getMembersInClubRequest struct {
	ClubId int `path:"clubId" minimum:"1"`
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
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	ok, err := h.authZService.IsMember(ctx, userID, req.ClubId)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to get members in this club")
	}

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
	MemberId int `param:"memberId" minimum:"1"`
}

func (h *Handler) RemoveMemberFromClub(ctx context.Context, req *removeUserFromClubRequest) (*struct{}, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	ok, err := h.authZService.IsAdmin(ctx, userID, req.MemberId)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to check authorization")
	}
	if !ok {
		return nil, huma.Error403Forbidden("user not authorized to remove member from club")
	}

	if err := h.memberService.DeleteMember(ctx, req.MemberId); err != nil {
		return nil, huma.Error500InternalServerError("failed to remove member from club, try again later")
	}

	return nil, nil
}
