package handlers

import (
	"context"
	"core/internal/authentication"
	"time"

	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type loginResponse struct {
	SetCookie []http.Cookie `header:"Set-Cookie"`
}

func (h *Handler) Login(ctx context.Context, req *loginRequest) (*loginResponse, error) {
	correct, accessToken, refreshToken, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to login, try again later")
	}
	if !correct {
		return nil, huma.Error400BadRequest("invalid email or password")
	}

	resp := &loginResponse{
		SetCookie: []http.Cookie{
			{
				Name:     "accessToken",
				Value:    accessToken,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				MaxAge:   int(authentication.AccessTokenDuration.Seconds()),
			},
			{
				Name:     "refreshToken",
				Value:    refreshToken,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				MaxAge:   int(authentication.RefreshTokenDuration.Seconds()),
			},
		},
	}

	return resp, nil
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type refreshResponse struct {
	SetCookie []http.Cookie `header:"Set-Cookie"`
}

// TODO: Could this just check cookies instead of requiring a request body?
func (h *Handler) Refresh(ctx context.Context, req *refreshRequest) (*refreshResponse, error) {
	valid, _, err := h.authService.VerifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to verify refresh token")
	}

	if !valid {
		return nil, huma.Error401Unauthorized("invalid refresh token")
	}

	accessToken, refreshToken, err := h.authService.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to refresh tokens")
	}

	resp := &refreshResponse{
		SetCookie: []http.Cookie{
			{
				Name:     "accessToken",
				Value:    accessToken,
				Path:     "/",
				Expires:  time.Now().Add(time.Hour),
				Secure:   true,
				HttpOnly: true,
			}, {
				Name:     "refreshToken",
				Value:    refreshToken,
				Path:     "/",
				Expires:  time.Now().Add(12 * time.Hour),
				Secure:   true,
				HttpOnly: true,
			},
		},
	}

	return resp, nil
}

type signupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=1,max=64"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (h *Handler) Signup(ctx context.Context, req *signupRequest) (*struct{}, error) {
	success, err := h.authService.Signup(ctx, req.Email, req.Name, req.Password)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to signup, try again later")
	}
	if !success {
		return nil, huma.Error400BadRequest("user with the given email already exists")
	}

	exists, _, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to signup, try again later")
	}
	if !exists {
		return nil, huma.Error500InternalServerError("failed to signup, try again later")
	}

	return nil, nil
}

type changePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

func (h *Handler) ChangePassword(ctx context.Context, req *changePasswordRequest) (*struct{}, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	if err := h.userService.UpdatePassword(ctx, userID, req.OldPassword, req.NewPassword); err != nil {
		return nil, huma.Error500InternalServerError("failed to change password, try again later")
	}

	return nil, nil
}

type logoutresponse struct {
	SetCookie []http.Cookie `header:"Set-Cookie"`
}

func (h *Handler) Logout(ctx context.Context, req *struct{}) (*logoutresponse, error) {
	resp := &logoutresponse{
		SetCookie: []http.Cookie{
			{
				Name:     "accessToken",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				MaxAge:   -1,
			},
			{
				Name:     "refreshToken",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				MaxAge:   -1,
			},
		},
	}

	return resp, nil
}
