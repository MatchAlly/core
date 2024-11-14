package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type updateUserRequest struct {
	Email string `json:"userId" format:"email"`
	Name  string `json:"name" minLength:"1" maxLength:"50"`
}

type updateUserResponse struct {
	Email string `json:"userId"`
	Name  string `json:"name"`
}

func (h *Handler) UpdateUser(ctx context.Context, req *updateUserRequest) (*updateUserResponse, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	if err := h.userService.UpdateUser(ctx, userID, req.Email, req.Name); err != nil {
		return nil, huma.Error500InternalServerError("failed to update user, try again later")
	}

	resp := &updateUserResponse{
		Email: req.Email,
		Name:  req.Name,
	}

	return resp, nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *struct{}) (*struct{}, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	if err := h.userService.DeleteUser(ctx, userID); err != nil {
		return nil, huma.Error500InternalServerError("failed to delete user, try again later")
	}

	return nil, nil
}
