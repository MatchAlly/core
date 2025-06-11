package club

import (
	"context"
	"core/internal/game"
	"core/internal/member"
	"core/internal/subscription"
	"fmt"
)

type Service interface {
	GetClub(ctx context.Context, id int) (*Club, error)
	GetClubs(ctx context.Context, ids []int) ([]Club, error)
	CreateClub(ctx context.Context, name string, userId int) (int, error)
	DeleteClub(ctx context.Context, id int) error
	UpdateClub(ctx context.Context, id int, name string) error
	GetGames(ctx context.Context, clubID int) ([]game.Game, error)
	CreateInvite(ctx context.Context, clubId, userId int, initiator Initiator) error
	GetPendingInvites(ctx context.Context, clubId int) ([]Invite, error)
	GetUserInvites(ctx context.Context, userId int) ([]Invite, error)
	AcceptInvite(ctx context.Context, inviteId int) error
	RejectInvite(ctx context.Context, inviteId int) error
}

type service struct {
	repo                Repository
	memberService       member.Service
	subscriptionService subscription.Service
}

func NewService(repo Repository, memberService member.Service, subscriptionService subscription.Service) Service {
	return &service{
		repo:                repo,
		memberService:       memberService,
		subscriptionService: subscriptionService,
	}
}

func (s *service) GetClub(ctx context.Context, id int) (*Club, error) {
	club, err := s.repo.GetClub(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get club: %w", err)
	}

	return club, nil
}

func (s *service) GetClubs(ctx context.Context, ids []int) ([]Club, error) {
	clubs, err := s.repo.GetClubs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get clubs: %w", err)
	}

	return clubs, nil
}

func (s *service) CreateClub(ctx context.Context, name string, userId int) (int, error) {
	// Validate club name
	if len(name) < 2 || len(name) > 50 {
		return 0, fmt.Errorf("club name must be between 2 and 50 characters")
	}

	// Create club
	club := &Club{
		Name: name,
	}
	clubId, err := s.repo.CreateClub(ctx, club)
	if err != nil {
		return 0, fmt.Errorf("failed to create club: %w", err)
	}

	// Create initial owner member
	member := &member.Member{
		ClubID: clubId,
		UserID: userId,
		Role:   member.RoleOwner,
	}
	if err := s.memberService.CreateMember(ctx, member); err != nil {
		// If member creation fails, delete the club to maintain consistency
		if delErr := s.repo.DeleteClub(ctx, clubId); delErr != nil {
			return 0, fmt.Errorf("failed to create member and cleanup club: %w (cleanup error: %v)", err, delErr)
		}
		return 0, fmt.Errorf("failed to create initial owner member: %w", err)
	}

	return clubId, nil
}

func (s *service) DeleteClub(ctx context.Context, id int) error {
	if err := s.repo.DeleteClub(ctx, id); err != nil {
		return fmt.Errorf("failed to delete club: %w", err)
	}

	return nil
}

func (s *service) UpdateClub(ctx context.Context, id int, name string) error {
	if len(name) < 2 || len(name) > 50 {
		return fmt.Errorf("club name must be between 2 and 50 characters")
	}

	if err := s.repo.UpdateClub(ctx, id, name); err != nil {
		return fmt.Errorf("failed to update club: %w", err)
	}

	return nil
}

func (s *service) GetGames(ctx context.Context, clubID int) ([]game.Game, error) {
	games, err := s.repo.GetGames(ctx, clubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get games: %w", err)
	}

	return games, nil
}

func (s *service) CreateInvite(ctx context.Context, clubId, userId int, initiator Initiator) error {
	// Check if user is already a member
	isMember, err := s.repo.IsMember(ctx, userId, clubId)
	if err != nil {
		return fmt.Errorf("failed to check membership: %w", err)
	}
	if isMember {
		return fmt.Errorf("user is already a member of this club")
	}

	// Check if invite already exists
	invites, err := s.repo.GetPendingInvites(ctx, clubId)
	if err != nil {
		return fmt.Errorf("failed to check existing invites: %w", err)
	}
	for _, invite := range invites {
		if invite.UserId == userId {
			return fmt.Errorf("user already has a pending invite")
		}
	}

	invite := &Invite{
		ClubId:    clubId,
		UserId:    userId,
		Initiator: initiator,
	}

	if err := s.repo.CreateInvite(ctx, invite); err != nil {
		return fmt.Errorf("failed to create invite: %w", err)
	}

	return nil
}

func (s *service) GetPendingInvites(ctx context.Context, clubId int) ([]Invite, error) {
	return s.repo.GetPendingInvites(ctx, clubId)
}

func (s *service) GetUserInvites(ctx context.Context, userId int) ([]Invite, error) {
	return s.repo.GetUserInvites(ctx, userId)
}

func (s *service) AcceptInvite(ctx context.Context, inviteId int) error {
	invite, err := s.repo.GetInvite(ctx, inviteId)
	if err != nil {
		return fmt.Errorf("failed to get invite: %w", err)
	}

	// Check if user is already a member
	isMember, err := s.repo.IsMember(ctx, invite.UserId, invite.ClubId)
	if err != nil {
		return fmt.Errorf("failed to check membership: %w", err)
	}
	if isMember {
		return fmt.Errorf("user is already a member of this club")
	}

	// Create member
	member := &member.Member{
		ClubID: invite.ClubId,
		UserID: invite.UserId,
		Role:   member.RoleMember,
	}
	if err := s.repo.CreateMember(ctx, member); err != nil {
		return fmt.Errorf("failed to create member: %w", err)
	}

	// Delete invite
	if err := s.repo.DeleteInvite(ctx, inviteId); err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}

	return nil
}

func (s *service) RejectInvite(ctx context.Context, inviteId int) error {
	if err := s.repo.DeleteInvite(ctx, inviteId); err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}

	return nil
}
