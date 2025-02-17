package user

import (
	"context"

	"encore.dev/beta/errs"
)

// LoginParams is the request data for the Login endpoint.
type LoginParams struct {
	Email    string
	Password string `encore:"sensitive"`
}

//encore:api public method=POST path=/login
func (s *Service) Login(ctx context.Context, params *LoginParams) error {
	var user User

	err := s.db.
		Where("email = $1", params.Email).
		First(&user).
		Error
	if err != nil {
		return errs.B().Code(errs.InvalidArgument).Msg("wrong email or password").Err()
	}

	passwordMatches := validatePassword(user.HashedPassword, params.Password)
	if !passwordMatches {
		return errs.B().Code(errs.InvalidArgument).Msg("wrong email or password").Err()
	}

	return nil
}
