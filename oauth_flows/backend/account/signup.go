package account

import (
	"context"
	"errors"

	"encore.app/internal/validation"
	"encore.app/oauth_flows/backend/account/db"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type SignupParams struct {
	Email                string `mod:"trim" validate:"required,email" encore:"sensitive" json:"email"`
	Password             string `mod:"trim" validate:"required,min=6,max=72" encore:"sensitive" json:"password"`
	PasswordConfirmation string `mod:"trim" validate:"required,eqcsfield=Password" encore:"sensitive" json:"password_confirmation"`
	// TODO: Get language from headers
}

type SignupResponse struct {
	ID int64 `json:"id"`
}

func (params *SignupParams) Validate() error {
	err := validation.ValidateTransform(context.Background(), params)
	return err
}

//encore:api public method=POST path=/signup
func (s *Service) Signup(ctx context.Context, params *SignupParams) (*SignupResponse, error) {
	hashedPassword, err := hashPassword(params.Password)
	if err != nil {
		rlog.Error("Error hashing User password.", "err", err)
		return nil, &errs.Error{Code: errs.Unknown, Message: "Unknown error"}
	}

	insertParams := db.InsertAccountParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	account, err := s.Query.InsertAccount(ctx, insertParams)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			details := &validation.ValidationErrors{
				"email": {"Email already taken"},
			}
			return nil, &errs.Error{Code: errs.InvalidArgument, Message: "", Details: details}
		}

		rlog.Error("Database error creating new User.", "err", err)
		return nil, &errs.Error{Code: errs.Unknown, Message: "Unknown error"}
	}

	_, err = EmailVerificationRequested.Publish(
		ctx,
		&EmailVerificationRequestedEvent{UserID: account.ID},
	)
	if err != nil {
		rlog.Warn("Error publishing EmailVerificationRequestedEvent.", "err", err)
	}

	return &SignupResponse{ID: account.ID}, nil
}

func hashPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)

	return string(hashedPasswordBytes), err
}

func validatePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		rlog.Error("error validationg password", "err", err)
	}
	return err == nil
}
