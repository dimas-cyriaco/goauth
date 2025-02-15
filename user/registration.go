package user

import (
	"context"
	"errors"

	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type RegistrationParams struct {
	Email                string `mod:"trim" validate:"required,email" encore:"sensitive"`
	Password             string `mod:"trim" validate:"required,min=6,max=72" encore:"sensitive"`
	PasswordConfirmation string `mod:"trim" validate:"required,eqcsfield=Password" encore:"sensitive"`
}

type RegistrationResponse struct {
	ID int `json:"id"`
}

func (params *RegistrationParams) Validate() error {
	validate := validator.New()
	conform := modifiers.New()

	err := conform.Struct(context.Background(), params)
	if err != nil {
		return err
	}

	return validate.Struct(params)
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
			return nil, &errs.Error{Code: errs.InvalidArgument, Message: "Invalid Argument"}
		}

		rlog.Error("Database error creating new User.", "err", err)
		return nil, &errs.Error{Code: errs.Unknown, Message: "Unknown error"}
	}

	_, err = Signups.Publish(ctx, &SignupEvent{UserID: user.ID})
	if err != nil {
		rlog.Error("Error publishing SignupEvent.", "err", err)
	}

	return &RegistrationResponse{ID: user.ID}, nil
}

func hashPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)

	return string(hashedPasswordBytes), err
}
