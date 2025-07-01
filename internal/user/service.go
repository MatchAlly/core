package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (exists bool, user *User, err error)
	CreateUser(ctx context.Context, email, name, hash string) (uuid.UUID, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateUser(ctx context.Context, id uuid.UUID, email, name string) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
}

type service struct {
	repo   Repository
	pepper string
}

func NewService(repo Repository, pepper string) Service {
	return &service{
		repo:   repo,
		pepper: pepper,
	}
}

func (s *service) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with id %s: %w", id, err)
	}

	return user, nil
}

func (s *service) CreateUser(ctx context.Context, email, name, hash string) (uuid.UUID, error) {
	user := &User{
		Email: email,
		Name:  name,
		Hash:  hash,
	}

	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (bool, *User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil, nil
		}

		return false, nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return true, user, nil
}

func (s *service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user with id %s: %w", id, err)
	}

	return nil
}

func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, email, name string) error {
	user := &User{
		ID:    id,
		Email: email,
		Name:  name,
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update user with id %s: %w", id, err)
	}

	return nil
}

func (s *service) UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	u, err := s.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user with id %s: %w", userID, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(oldPassword+s.pepper)); err != nil {
		return fmt.Errorf("failed to compare password: %w", err)
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword+s.pepper), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.repo.UpdatePassword(ctx, userID, string(hashedPasswordBytes)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *service) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.UpdateLastLogin(ctx, userID); err != nil {
		return fmt.Errorf("failed to update last login for user with id %s: %w", userID, err)
	}

	return nil
}
