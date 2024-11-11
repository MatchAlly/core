package authentication

import (
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 24 * time.Hour
)

type AccessClaims struct {
	jwt.StandardClaims

	Name string `json:"name"`
}

type RefreshClaims struct {
	jwt.StandardClaims
}
