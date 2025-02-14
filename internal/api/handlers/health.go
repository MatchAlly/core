package handlers

import (
	"context"
)

type healthRequest struct{}

type healthResponse struct {
	Status string `json:"status"`
}

func (h *Handler) Health(ctx context.Context, req *healthRequest) (*healthResponse, error) {
	return &healthResponse{
		Status: "healthy",
	}, nil
}
