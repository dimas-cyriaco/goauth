package user

import (
	"context"

	"encore.dev/rlog"
)

type RequestVerificationEmailParams struct {
	Email string
}

//encore:api public method=POST path=/request_email_verification
func (s *Service) RequestVerificationEmail(ctx context.Context, params *RequestVerificationEmailParams) error {
	var user User
	if err := s.db.Where("email = $1", params.Email).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return nil
		}

		rlog.Error("Error reading User from DB.", "err", err)

		return err
	}

	if user.EmailVerifiedAt != nil {
		return nil
	}

	_, err := EmailVerificationRequested.Publish(
		ctx,
		&EmailVerificationRequestedEvent{UserID: user.ID},
	)
	if err != nil {
		rlog.Error("Error publishing EmailVerificationRequestedEvent.", "err", err)
		return err
	}

	return nil
}
