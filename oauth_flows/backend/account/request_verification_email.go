package account

import (
	"context"

	"encore.dev/rlog"
)

type RequestVerificationEmailParams struct {
	Email string
}

//encore:api public method=POST path=/request_email_verification
func (s *Service) RequestVerificationEmail(ctx context.Context, params *RequestVerificationEmailParams) error {
	account, err := s.Query.FindAccountByEmail(ctx, params.Email)
	if err != nil {
		if err.Error() == "record not found" {
			return nil
		}

		rlog.Error("Error reading Account from DB.", "err", err)

		return err
	}

	if account.EmailVerifiedAt.Valid {
		return nil
	}

	_, err = EmailVerificationRequested.Publish(
		ctx,
		&EmailVerificationRequestedEvent{AccountID: account.ID},
	)
	if err != nil {
		rlog.Error("Error publishing EmailVerificationRequestedEvent.", "err", err)
		return err
	}

	return nil
}
