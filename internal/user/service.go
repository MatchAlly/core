package user

import (
	"context"

	"github.com/pkg/errors"
)

type Service interface {
	GetUser(ctx context.Context, id uint) (*User, error)
	GetUsers(ctx context.Context, ids []uint) ([]User, error)
	GetUserByEmail(ctx context.Context, email string) (exists bool, user *User, err error)
	GetUsersByEmails(ctx context.Context, emails []string) ([]User, error)
	CreateUser(ctx context.Context, email, name, hash string) error
	DeleteUser(ctx context.Context, id uint) error
	UpdateUser(ctx context.Context, id uint, email, name string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetUser(ctx context.Context, id uint) (*User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get user %d", id)
	}

	return user, nil
}

func (s *service) GetUsers(ctx context.Context, ids []uint) ([]User, error) {
	users, err := s.repo.GetUsers(ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get users %v", ids)
	}

	return users, nil
}

func (s *service) CreateUser(ctx context.Context, email, name, hash string) error {
	user := &User{
		Email: email,
		Name:  name,
		Hash:  hash,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
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

func (s *service) GetUsersByEmails(ctx context.Context, emails []string) ([]User, error) {
	users, err := s.repo.GetUsersByEmails(ctx, emails)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users by emails")
	}

	return users, nil
}

func (s *service) DeleteUser(ctx context.Context, id uint) error {
	if err := s.repo.DeleteUser(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete user by id")
	}

	return nil
}

func (s *service) UpdateUser(ctx context.Context, id uint, email, name string) error {
	user := &User{
		Id:    id,
		Email: email,
		Name:  name,
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}
