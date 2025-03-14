package account

import (
	"context"
	"strconv"

	"encore.app/oauth_flows/backend/internal/tokens"
	"encore.dev/rlog"
)

type VerifyEmailParams struct {
	Token string `query:"verification_token"`
}

//encore:api public method=GET path=/verify_email
func (s *Service) VerifyEmail(ctx context.Context, params *VerifyEmailParams) error {
	payload, err := tokens.GetPayloadForToken(tokens.EmailVerification, params.Token)
	if err != nil {
		return err
	}

	accountID, err := strconv.ParseInt(payload["AccountID"], 10, 64)
	if err != nil {
		rlog.Error("Error converting accountID to int64.", "err", err)
		return nil
	}

	s.Query.MarkEmailAsVerified(ctx, accountID)

	return nil
}
