package handlers

import (
	"context"
	"core/internal/subscription"
	"fmt"
	"time"

	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type loginRequest struct {
	Body struct {
		Email    string `json:"email" format:"email"`
		Password string `json:"password" minLength:"8" maxLength:"256"`
	}
}

type loginResponse struct {
	SetCookie []http.Cookie `header:"Set-Cookie"`
}

func (h *Handler) Login(ctx context.Context, req *loginRequest) (*loginResponse, error) {
	correct, accessToken, refreshToken, err := h.authNService.Login(ctx, req.Body.Email, req.Body.Password)
	if err != nil {
		h.l.Errorw("failed to login", "error", err)
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
				MaxAge:   int(h.config.AccessTokenDuration.Seconds()),
			},
			{
				Name:     "refreshToken",
				Value:    refreshToken,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				MaxAge:   int(h.config.RefreshTokenDuration.Seconds()),
			},
		},
	}

	return resp, nil
}

type refreshRequest struct {
	RefreshToken http.Cookie `json:"refreshToken"`
}

type refreshResponse struct {
	SetCookie []http.Cookie `header:"Set-Cookie"`
}

func (h *Handler) Refresh(ctx context.Context, req *refreshRequest) (*refreshResponse, error) {
	valid, _, err := h.authNService.VerifyRefreshToken(ctx, req.RefreshToken.Value)
	if err != nil {
		h.l.Error("failed to verify refresh token", "error", err)
		return nil, huma.Error500InternalServerError("failed to verify refresh token")
	}

	if !valid {
		return nil, huma.Error401Unauthorized("invalid refresh token")
	}

	accessToken, refreshToken, err := h.authNService.RefreshTokens(ctx, req.RefreshToken.Value)
	if err != nil {
		h.l.Error("failed to refresh tokens", "error", err)
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
	Body struct {
		Email    string `json:"email" format:"email"`
		Name     string `json:"name" minLength:"2" maxLength:"20"`
		Password string `json:"password" minLength:"8" maxLength:"256"`
	}
}

func (h *Handler) Signup(ctx context.Context, req *signupRequest) (*struct{}, error) {
	fmt.Println("Signup request:", req)
	success, err := h.authNService.Signup(ctx, req.Body.Email, req.Body.Name, req.Body.Password)
	if err != nil {
		h.l.Errorw("failed to signup", "error", err)
		return nil, huma.Error500InternalServerError("failed to signup, try again later")
	}
	if !success {
		return nil, huma.Error400BadRequest("user with the given email already exists")
	}

	exists, u, err := h.userService.GetUserByEmail(ctx, req.Body.Email)
	if err != nil || !exists {
		h.l.Errorw("failed to get user by email", "error", err)
		return nil, huma.Error500InternalServerError("failed to signup, try again later")
	}

	sub := subscription.Subscription{
		UserID:    u.ID,
		Tier:      subscription.TierFree,
		CreatedAt: time.Now(),
	}
	if err := h.subscriptionService.CreateSubscription(ctx, sub); err != nil {
		h.l.Errorw("failed to create subscription", "error", err)
		return nil, huma.Error500InternalServerError("failed to create subscription")
	}

	return nil, nil
}

type changePasswordRequest struct {
	Body struct {
		OldPassword string `json:"oldPassword" minLength:"8" maxLength:"256"`
		NewPassword string `json:"newPassword" minLength:"8" maxLength:"256"`
	}
}

func (h *Handler) ChangePassword(ctx context.Context, req *changePasswordRequest) (*struct{}, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	if err := h.userService.UpdatePassword(ctx, userID, req.Body.OldPassword, req.Body.NewPassword); err != nil {
		h.l.Errorw("failed to change password", "error", err)
		return nil, huma.Error500InternalServerError("failed to change password, try again later")
	}

	return nil, nil
}

type logoutResponse struct {
	SetCookie []http.Cookie `header:"Set-Cookie"`
}

func (h *Handler) Logout(ctx context.Context, req *struct{}) (*logoutResponse, error) {
	token, ok := ctx.Value("token").(string)
	if !ok {
		h.l.Error("failed to get token from context")
		return nil, huma.Error500InternalServerError("failed to get token from context")
	}

	if err := h.authNService.Logout(ctx, token); err != nil {
		h.l.Errorw("failed to logout", "error", err)
		return nil, huma.Error500InternalServerError("failed to logout, try again later")
	}

	resp := &logoutResponse{
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
