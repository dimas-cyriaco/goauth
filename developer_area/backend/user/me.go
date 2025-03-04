package user

import (
	"context"
	"time"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
)

type MeResponse struct {
	ID              int        `json:"id"`
	Email           string     `json:"email" encore:"sensitive"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" faker:"-"`
}

//encore:api auth method=GET path=/me
func (s *Service) Me(ctx context.Context) (*MeResponse, error) {
	userID, hasUserId := auth.UserID()
	if !hasUserId {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "unauthenticated"}
	}

	var user User
	if err := s.db.Where("id = $1", userID).First(&user).Error; err != nil {
		return nil, err
	}

	response := MeResponse{
		ID:              user.ID,
		Email:           user.Email,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		EmailVerifiedAt: user.EmailVerifiedAt,
	}
	return &response, nil
}
