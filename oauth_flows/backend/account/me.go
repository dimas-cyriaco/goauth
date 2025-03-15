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
	accountID, hasAccountId := auth.UserID()
	if !hasAccountId {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "unauthenticated"}
	}

	aID, err := strconv.ParseInt(string(accountID), 10, 64)
	if err != nil {
		rlog.Error("Error converting sessionID to int64.", "err", err)
		return nil, err
	}

	account, err := s.Query.FindAccountByID(ctx, aID)
	if err != nil {
		return nil, err
	}

	response := MeResponse{
		ID:              account.ID,
		Email:           account.Email,
		CreatedAt:       account.CreatedAt.Time,
		UpdatedAt:       account.UpdatedAt.Time,
		EmailVerifiedAt: &account.EmailVerifiedAt.Time,
	}
	return &response, nil
}
