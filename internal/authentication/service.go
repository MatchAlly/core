//go:generate mockgen --source=service.go -destination=service_mock.go -package=authentication
package authentication

import (
	"context"
	"core/internal/user"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Secret        string        `mapstructure:"secret" validate:"required"`
	AccessExpiry  time.Duration `mapstructure:"access_expiry" validate:"required"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry" validate:"required"`
}

type Service interface {
	Login(ctx context.Context, email, password string) (valid bool, accessToken, refreshToken string, err error)
	Signup(ctx context.Context, email, username, password string) (success bool, err error)
	VerifyAccessToken(ctx context.Context, token string) (valid bool, claims *AccessClaims, err error)
	VerifyRefreshToken(ctx context.Context, token string) (valid bool, claims *RefreshClaims, err error)
	RefreshTokens(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error)
}

type service struct {
	config      Config
	userService user.Service
}

func NewService(config Config, userService user.Service) Service {
	return &service{
		config:      config,
		userService: userService,
	}
}

func (s *service) Login(ctx context.Context, email string, password string) (bool, string, string, error) {
	exists, user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return false, "", "", errors.Wrap(err, "failed to get user by email")
	}

	if !exists {
		return false, "", "", nil
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password)); err != nil {
		return false, "", "", nil
	}

	accessToken, refreshToken, err := s.generateTokenPair(user.Name, user.Id)
	if err != nil {
		return false, "", "", errors.Wrap(err, "failed to generate jwts")
	}

	return true, accessToken, refreshToken, nil
}

func (s *service) Signup(ctx context.Context, email string, username string, password string) (bool, error) {
	exists, _, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return false, errors.Wrap(err, "failed to check for existing user with email")
	}

	if exists {
		return false, nil
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return false, errors.Wrap(err, "failed to hash password")
	}

	if err = s.userService.CreateUser(ctx, email, username, string(hashedPasswordBytes)); err != nil {
		return false, errors.Wrap(err, "failed to create user")
	}

	return true, nil
}

func (s *service) VerifyAccessToken(ctx context.Context, token string) (bool, *AccessClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Secret), nil
	})
	if err != nil {
		return false, nil, err
	}

	if _, ok := parsedToken.Method.(*jwt.SigningMethodHMAC); !ok {
		return false, nil, errors.New("unexpected signing method")
	}

	if !parsedToken.Valid {
		return false, nil, errors.New("jwt is invalid")
	}

	claims, ok := parsedToken.Claims.(*AccessClaims)
	if !ok {
		return false, nil, errors.New("failed to parse claims")
	}

	return true, claims, nil
}

func (s *service) VerifyRefreshToken(ctx context.Context, token string) (bool, *RefreshClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Secret), nil
	})
	if err != nil {
		return false, nil, err
	}

	if _, ok := parsedToken.Method.(*jwt.SigningMethodHMAC); !ok {
		return false, nil, errors.New("unexpected signing method")
	}

	if !parsedToken.Valid {
		return false, nil, errors.New("jwt is invalid")
	}

	claims, ok := parsedToken.Claims.(*RefreshClaims)
	if !ok {
		return false, nil, errors.New("failed to parse claims")
	}

	return true, claims, nil
}

func (s *service) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	_, claims, err := s.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to verify refresh token")
	}

	userId, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to parse user id")
	}

	u, err := s.userService.GetUser(ctx, uint(userId))
	if err != nil {
		return "", "", errors.Wrap(err, "failed to get user by id")
	}

	accessToken, newRefreshToken, err := s.generateTokenPair(u.Name, uint(userId))
	if err != nil {
		return "", "", errors.Wrap(err, "failed to generate jwts")
	}

	return accessToken, newRefreshToken, nil
}

func (s *service) generateTokenPair(name string, userId uint) (string, string, error) {
	now := time.Now()

	userIdString := strconv.FormatUint(uint64(userId), 10)

	accessclaims := AccessClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   userIdString,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(15 * time.Minute).Unix(),
			Issuer:    "MatchAlly",
		},
		Name: name,
	}

	accessTokenUnsigned := jwt.NewWithClaims(jwt.SigningMethodHS256, accessclaims)
	accessTokenSigned, err := accessTokenUnsigned.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", "", errors.Wrap(err, "failed to sign access token")
	}

	refreshClaims := RefreshClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   userIdString,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(12 * time.Hour).Unix(),
			Issuer:    "MatchAlly",
		},
	}

	refreshTokenUnsigned := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenSigned, err := refreshTokenUnsigned.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", "", errors.Wrap(err, "failed to sign refresh token")
	}

	return accessTokenSigned, refreshTokenSigned, nil
}
