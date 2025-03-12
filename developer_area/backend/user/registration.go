package user

import (
	"context"
	"errors"

	"encore.app/developer_area/backend/internal/utils"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type RegistrationParams struct {
	Email                string `mod:"trim" validate:"required,email" encore:"sensitive" json:"email"`
	Password             string `mod:"trim" validate:"required,min=6,max=72" encore:"sensitive" json:"password"`
	PasswordConfirmation string `mod:"trim" validate:"required,eqcsfield=Password" encore:"sensitive" json:"password_confirmation"`
	// TODO: Get language from headers
}

type RegistrationResponse struct {
	ID int `json:"id"`
}

func (params *RegistrationParams) Validate() error {
	err := utils.ValidateTransform(context.Background(), params)
	return err
}

//encore:api public method=POST path=/sign_up
func (s *Service) Registration(ctx context.Context, params *RegistrationParams) (*RegistrationResponse, error) {
	hashedPassword, err := hashPassword(params.Password)
	if err != nil {
		rlog.Error("Error hashing User password.", "err", err)
		return nil, &errs.Error{Code: errs.Unknown, Message: "Unknown error"}
	}

	user := &User{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	err = s.db.Create(user).Error
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			details := &utils.ValidationErrors{
				"email": {"Email already taken"},
			}
			return nil, &errs.Error{Code: errs.InvalidArgument, Message: "", Details: details}
		}

		rlog.Error("Database error creating new User.", "err", err)
		return nil, &errs.Error{Code: errs.Unknown, Message: "Unknown error"}
	}

	_, err = EmailVerificationRequested.Publish(
		ctx,
		&EmailVerificationRequestedEvent{UserID: user.ID},
	)
	if err != nil {
		rlog.Warn("Error publishing EmailVerificationRequestedEvent.", "err", err)
	}

	return &RegistrationResponse{ID: user.ID}, nil
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
