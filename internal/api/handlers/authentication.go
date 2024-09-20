package handlers

import (
	"core/internal/api/helpers"
	"time"

	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (h *Handler) Login(c echo.Context) error {
	req, ctx, err := helpers.Bind[loginRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	correct, accessToken, refreshToken, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.l.Error("failed to login", "error", err)
		return echo.ErrInternalServerError
	}
	if !correct {
		return echo.ErrUnauthorized
	}

	c.SetCookie(&http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		Secure:   true,
		HttpOnly: true,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(12 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	})

	return c.NoContent(http.StatusOK)
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

func (h *Handler) Refresh(c echo.Context) error {
	req, ctx, err := helpers.Bind[refreshRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	valid, _, err := h.authService.VerifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		h.l.Error("failed to verify refresh token", "error", err)
		return echo.ErrInternalServerError
	}

	if !valid {
		return echo.ErrUnauthorized
	}

	accessToken, refreshToken, err := h.authService.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		h.l.Error("failed to generate new tokens", "error", err)
		return echo.ErrInternalServerError
	}

	c.SetCookie(&http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		Secure:   true,
		HttpOnly: true,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(12 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	})

	return c.NoContent(http.StatusOK)
}

type signupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=1,max=64"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (h *Handler) Signup(c echo.Context) error {
	req, ctx, err := helpers.Bind[signupRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	success, err := h.authService.Signup(ctx, req.Email, req.Name, req.Password)
	if err != nil {
		h.l.Error("failed to signup", zap.Error(err))
		return echo.ErrInternalServerError
	}
	if !success {
		h.l.Debug("cant sign up, user already exists", zap.String("name", req.Name))
		return echo.ErrBadRequest
	}

	exists, _, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		h.l.Error("failed to get user by email", zap.Error(err))
		return echo.ErrInternalServerError
	}
	if !exists {
		h.l.Debug("signed up user not found", zap.String("name", req.Name))
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusCreated)
}

type changePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

func (h *Handler) ChangePassword(c helpers.AuthContext) error {
	req, ctx, err := helpers.Bind[changePasswordRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err := h.userService.UpdatePassword(ctx, c.UserID, req.OldPassword, req.NewPassword); err != nil {
		h.l.Error("failed to delete user", "error", err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
