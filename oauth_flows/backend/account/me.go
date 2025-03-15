package account

import (
	"context"
	"strconv"
	"time"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

type MeResponse struct {
	ID              int64      `json:"id"`
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

	uid, err := strconv.ParseInt(string(userID), 10, 64)
	if err != nil {
		rlog.Error("Error converting sessionID to int64.", "err", err)
		return nil, err
	}

	// var user User
	user, err := s.Query.FindAccountByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	response := MeResponse{
		ID:              user.ID,
		Email:           user.Email,
		CreatedAt:       user.CreatedAt.Time,
		UpdatedAt:       user.UpdatedAt.Time,
		EmailVerifiedAt: &user.EmailVerifiedAt.Time,
	}
	return &response, nil
}
