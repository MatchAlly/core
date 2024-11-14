package user

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetUser(ctx context.Context, id int) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (exists bool, user *User, err error)
	CreateUser(ctx context.Context, email, name, hash string) (int, error)
	DeleteUser(ctx context.Context, id int) error
	UpdateUser(ctx context.Context, id int, email, name string) error
	UpdatePassword(ctx context.Context, userID int, oldPassword, newPassword string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetUser(ctx context.Context, id int) (*User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get user %d", id)
	}

	return user, nil
}

func (s *service) CreateUser(ctx context.Context, email, name, hash string) (int, error) {
	user := &User{
		Email: email,
		Name:  name,
		Hash:  hash,
	}

	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create user")
	}

	return id, nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (bool, *User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil, nil
		}

		return false, nil, errors.Wrap(err, "failed to get user by email")
	}

	return true, user, nil
}

func (s *service) DeleteUser(ctx context.Context, id int) error {
	if err := s.repo.DeleteUser(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete user by id")
	}

	return nil
}

func (s *service) UpdateUser(ctx context.Context, id int, email, name string) error {
	user := &User{
		ID:    id,
		Email: email,
		Name:  name,
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}

func (s *service) UpdatePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	u, err := s.GetUser(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "failed to get user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(oldPassword)); err != nil {
		return errors.Wrap(err, "incorrect password")
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.MinCost)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	if err := s.repo.UpdatePassword(ctx, userID, string(hashedPasswordBytes)); err != nil {
		return errors.Wrap(err, "failed to update password")
	}

	return nil
}
