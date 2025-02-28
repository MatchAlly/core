package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service interface {
	SetTokenUsed(ctx context.Context, token string) error
	GetTokenUsed(ctx context.Context, token string) (bool, error)
}

type service struct {
	client         *redis.Client
	denylistExpiry time.Duration
}

func NewService(client *redis.Client, denylistExpiry time.Duration) Service {
	return &service{
		denylistExpiry: denylistExpiry,
		client:         client,
	}
}

func (s *service) SetTokenUsed(ctx context.Context, token string) error {
	key := denylistTokenKey(token)
	return s.client.Set(ctx, key, true, s.denylistExpiry).Err()
}

func (s *service) GetTokenUsed(ctx context.Context, token string) (bool, error) {
	key := denylistTokenKey(token)
	return s.client.Get(ctx, key).Bool()
}
