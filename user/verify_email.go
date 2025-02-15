package user

import (
	"context"
	"time"

	tokengenerator "encore.app/internal/token_generator"
)

type VerifyEmailParams struct {
	Token string `query:"verification_token"`
}

//encore:api public method=GET path=/verify_email
func (s *Service) VerifyEmail(ctx context.Context, params *VerifyEmailParams) error {
	payload, err := tokengenerator.GetPayloadForToken(tokengenerator.EmailVerification, params.Token)
	if err != nil {
		return err
	}

	s.db.
		Model(&User{}).
		Where("ID = ?", payload["UserID"]).
		Where("email_verified_at is NULL").
		Update("email_verified_at", time.Now())

	return nil
}
