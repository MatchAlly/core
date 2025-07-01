package authorization

import (
	"context"
	"core/internal/member"

	"github.com/google/uuid"
)

type Service interface {
	IsMember(ctx context.Context, userID, clubID uuid.UUID) (bool, error)
	IsAdmin(ctx context.Context, userID, clubID uuid.UUID) (bool, error)
}

type service struct {
	memberService member.Service
}

func NewService(memberService member.Service) Service {
	return &service{
		memberService: memberService,
	}
}

func (s *service) IsMember(ctx context.Context, userID, clubID uuid.UUID) (bool, error) {
	memberships, err := s.memberService.GetUserMemberships(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, membership := range memberships {
		if membership.UserID == userID {
			return true, nil
		}
	}

	return false, nil
}

func (s *service) IsAdmin(ctx context.Context, userID, clubID uuid.UUID) (bool, error) {
	memberships, err := s.memberService.GetUserMemberships(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, membership := range memberships {
		if membership.UserID == userID && membership.Role == member.RoleAdmin {
			return true, nil
		}
	}

	return false, nil
}
