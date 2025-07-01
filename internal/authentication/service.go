package authentication

import (
	"context"
	"core/internal/cache"
	"core/internal/subscription"
	"core/internal/user"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Secret        string        `mapstructure:"secret" validate:"required"`
	AccessExpiry  time.Duration `mapstructure:"access_expiry" validate:"required"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry" validate:"required"`
	Pepper        string        `mapstructure:"pepper" validate:"required"`
}

type Service interface {
	Login(ctx context.Context, email, password string) (valid bool, accessToken, refreshToken string, err error)
	Logout(ctx context.Context, token string) error
	Signup(ctx context.Context, email, username, password string) (success bool, err error)
	VerifyAccessToken(ctx context.Context, token string) (valid bool, claims *AccessClaims, err error)
	VerifyRefreshToken(ctx context.Context, token string) (valid bool, claims *RefreshClaims, err error)
	RefreshTokens(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error)
}

type service struct {
	config              Config
	userService         user.Service
	subscriptionService subscription.Service
	cache               cache.Service
}

func NewService(config Config, userService user.Service, subscriptionService subscription.Service, cache cache.Service) Service {
	return &service{
		config:              config,
		userService:         userService,
		subscriptionService: subscriptionService,
		cache:               cache,
	}
}

func (s *service) Login(ctx context.Context, email string, password string) (bool, string, string, error) {
	exists, user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return false, "", "", fmt.Errorf("failed to check for existing user with email: %w", err)
	}
	if !exists {
		return false, "", "", nil
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password+s.config.Pepper)); err != nil {
		return false, "", "", nil
	}

	if err := s.userService.UpdateLastLogin(ctx, user.ID); err != nil {
		return false, "", "", fmt.Errorf("failed to update last login: %w", err)
	}

	accessToken, refreshToken, err := s.generateTokenPair(ctx, user.ID)
	if err != nil {
		return false, "", "", fmt.Errorf("failed to generate jwts: %w", err)
	}

	return true, accessToken, refreshToken, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	if err := s.cache.SetTokenUsed(ctx, token); err != nil {
		return fmt.Errorf("failed to set token as used: %w", err)

	}

	return nil
}

func (s *service) Signup(ctx context.Context, email string, username string, password string) (bool, error) {
	exists, _, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check for existing user with email: %w", err)
	}
	if exists {
		return false, nil
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password+s.config.Pepper), bcrypt.MinCost)
	if err != nil {
		return false, fmt.Errorf("failed to hash password: %w", err)
	}

	if _, err = s.userService.CreateUser(ctx, email, username, string(hashedPasswordBytes)); err != nil {
		return false, fmt.Errorf("failed to create user: %w", err)
	}

	return true, nil
}

func (s *service) VerifyAccessToken(ctx context.Context, token string) (bool, *AccessClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &AccessClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})
	if err != nil {
		return false, nil, err
	}
	if !parsedToken.Valid {
		return false, nil, fmt.Errorf("jwt is invalid")
	}

	claims, ok := parsedToken.Claims.(*AccessClaims)
	if !ok {
		return false, nil, fmt.Errorf("failed to parse claims")
	}

	return true, claims, nil
}

func (s *service) VerifyRefreshToken(ctx context.Context, token string) (bool, *RefreshClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &RefreshClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})
	if err != nil {
		return false, nil, err
	}
	if !parsedToken.Valid {
		return false, nil, fmt.Errorf("jwt is invalid")
	}

	claims, ok := parsedToken.Claims.(*RefreshClaims)
	if !ok {
		return false, nil, fmt.Errorf("failed to parse claims")
	}

	return true, claims, nil
}

func (s *service) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	_, claims, err := s.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to verify refresh token: %w", err)
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse user id: %w", err)
	}

	accessToken, newRefreshToken, err := s.generateTokenPair(ctx, userId)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token pair: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

func (s *service) generateTokenPair(ctx context.Context, userId uuid.UUID) (string, string, error) {
	sub, err := s.subscriptionService.GetByUserID(ctx, userId)
	if err != nil {
		return "", "", fmt.Errorf("failed to get subscription: %w", err)
	}

	now := time.Now()

	accessclaims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.AccessExpiry)),
		},
		SubscriptionTier: sub.Tier,
	}

	accessTokenUnsigned := jwt.NewWithClaims(jwt.SigningMethodHS256, accessclaims)
	accessTokenSigned, err := accessTokenUnsigned.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshClaims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.RefreshExpiry)),
		},
	}

	refreshTokenUnsigned := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenSigned, err := refreshTokenUnsigned.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return accessTokenSigned, refreshTokenSigned, nil
}
