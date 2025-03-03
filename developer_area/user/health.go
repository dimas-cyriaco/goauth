package user

import (
	"context"
)

type HealthResponse struct {
	Success bool `json:"success"`
}

//encore:api public method=GET path=/health
func (s *Service) Health(ctx context.Context) (*HealthResponse, error) {
	return &HealthResponse{Success: true}, nil
}
