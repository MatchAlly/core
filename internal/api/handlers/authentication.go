package handlers

import (
	"core/internal/api/helpers"

	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type loginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[loginRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	correct, accessToken, refreshToken, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.l.Error("failed to login",
			"error", err)
		return echo.ErrInternalServerError
	}
	if !correct {
		return echo.ErrUnauthorized
	}

	resp := loginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.JSON(http.StatusOK, resp)
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type refreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) Refresh(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[refreshRequest](c)
	if err != nil {
		return echo.ErrBadRequest
	}

	valid, _, err := h.authService.VerifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		h.l.Error("failed to verify refresh token",
			"error", err)
		return echo.ErrInternalServerError
	}

	if !valid {
		return echo.ErrUnauthorized
	}

	accessToken, refreshToken, err := h.authService.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		h.l.Error("failed to generate new tokens",
			"error", err)
		return echo.ErrInternalServerError
	}

	resp := refreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.JSON(http.StatusOK, resp)
}

type signupRequest struct {
	Email    string `json:"email" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *Handler) Signup(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := helpers.Bind[signupRequest](c)
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
