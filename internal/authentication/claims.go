package authentication

import (
	"github.com/golang-jwt/jwt"
)

type AccessClaims struct {
	jwt.StandardClaims

	Name string `json:"name"`
}

type RefreshClaims struct {
	jwt.StandardClaims
}
