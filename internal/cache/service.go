package cache

import (
	"context"
	"time"

	"github.com/valkey-io/valkey-go"
)

type Service interface {
	SetTokenUsed(ctx context.Context, token string) error
	GetTokenUsed(ctx context.Context, token string) (bool, error)
}

type service struct {
	client         valkey.Client
	denylistExpiry time.Duration
}

func NewService(client valkey.Client, denylistExpiry time.Duration) Service {
	return &service{
		denylistExpiry: denylistExpiry,
		client:         client,
	}
}

func (s *service) SetTokenUsed(ctx context.Context, token string) error {
	key := s.denylistTokenKey(token)
	return s.client.Do(ctx, s.client.B().Set().Key(key).Value("1").Nx().Ex(s.denylistExpiry).Build()).Error()
}

func (s *service) GetTokenUsed(ctx context.Context, token string) (bool, error) {
	key := s.denylistTokenKey(token)
	return s.client.Do(ctx, s.client.B().Get().Key(key).Build()).AsBool()
}

func (s *service) denylistTokenKey(tokenID string) string {
	return "denylist:" + tokenID
}
