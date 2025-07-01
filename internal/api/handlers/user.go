package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type updateUserRequest struct {
	UserID uuid.UUID `path:"userId"`
	Body   struct {
		Email string `json:"userId" format:"email"`
		Name  string `json:"name" minLength:"1" maxLength:"50"`
	}
}

type updateUserResponse struct {
	Body struct {
		Email string `json:"userId"`
		Name  string `json:"name"`
	}
}

func (h *Handler) UpdateUser(ctx context.Context, req *updateUserRequest) (*updateUserResponse, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	if userID != req.UserID {
		h.l.Error("user id from context does not match request")
		return nil, huma.Error403Forbidden("user id from context does not match request")
	}

	if err := h.user.UpdateUser(ctx, userID, req.Body.Email, req.Body.Name); err != nil {
		h.l.Error("failed to update user", "error", err)
		return nil, huma.Error500InternalServerError("failed to update user, try again later")
	}

	resp := &updateUserResponse{}
	resp.Body.Email = req.Body.Email
	resp.Body.Name = req.Body.Name

	return resp, nil
}

type deleteUserRequest struct {
	UserID uuid.UUID `path:"userId"`
}

func (h *Handler) DeleteUser(ctx context.Context, req *deleteUserRequest) (*struct{}, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		h.l.Error("failed to get user id from context")
		return nil, huma.Error500InternalServerError("failed to get user id from context")
	}

	if userID != req.UserID {
		h.l.Error("user id from context does not match request")
		return nil, huma.Error403Forbidden("user id from context does not match request")
	}

	if err := h.user.DeleteUser(ctx, userID); err != nil {
		h.l.Error("failed to delete user", "error", err)
		return nil, huma.Error500InternalServerError("failed to delete user, try again later")
	}

	return nil, nil
}
