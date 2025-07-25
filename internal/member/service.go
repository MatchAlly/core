package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	GetMembersInClub(ctx context.Context, clubId uuid.UUID) ([]Member, error)
	GetUserMemberships(ctx context.Context, userId uuid.UUID) ([]Member, error)
	GetMember(ctx context.Context, id uuid.UUID) (*Member, error)
	CreateMember(ctx context.Context, member *Member) error
	UpdateRole(ctx context.Context, memberId uuid.UUID, role Role) error
	DeleteMember(ctx context.Context, memberId uuid.UUID) error
	IsMember(ctx context.Context, userId, clubId uuid.UUID) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetMembersInClub(ctx context.Context, clubId uuid.UUID) ([]Member, error) {
	return s.repo.GetMembersInClub(ctx, clubId)
}

func (s *service) GetUserMemberships(ctx context.Context, userId uuid.UUID) ([]Member, error) {
	return s.repo.GetUserMemberships(ctx, userId)
}

func (s *service) GetMember(ctx context.Context, id uuid.UUID) (*Member, error) {
	return s.repo.GetMember(ctx, id)
}

func (s *service) CreateMember(ctx context.Context, member *Member) error {
	if err := s.repo.CreateMember(ctx, member); err != nil {
		return fmt.Errorf("failed to create member: %w", err)
	}
	return nil
}

func (s *service) UpdateRole(ctx context.Context, memberId uuid.UUID, role Role) error {
	if err := s.repo.UpdateRole(ctx, memberId, role); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

func (s *service) DeleteMember(ctx context.Context, memberId uuid.UUID) error {
	if err := s.repo.DeleteMember(ctx, memberId); err != nil {
		return fmt.Errorf("failed to delete member: %w", err)
	}
	return nil
}

func (s *service) IsMember(ctx context.Context, userId, clubId uuid.UUID) (bool, error) {
	memberships, err := s.repo.GetUserMemberships(ctx, userId)
	if err != nil {
		return false, err
	}

	for _, m := range memberships {
		if m.ClubID == clubId {
			return true, nil
		}
	}

	return false, nil
}
