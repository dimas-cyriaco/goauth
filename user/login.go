package user

import (
	"context"
	"strconv"

	tokengenerator "encore.app/internal/token_generator"
	"encore.dev/beta/errs"
)

// LoginParams is the request data for the Login endpoint.
type LoginParams struct {
	Email    string
	Password string `encore:"sensitive"`
}

type LoginResponse struct {
	SessionToken string `header:"Set-Cookie"`
}

//encore:api public method=POST path=/login
func (s *Service) Login(ctx context.Context, params *LoginParams) (*LoginResponse, error) {
	var user User

	err := s.db.
		Where("email = $1", params.Email).
		First(&user).
		Error
	if err != nil {
		return nil, errs.B().Code(errs.InvalidArgument).Msg("wrong email or password").Err()
	}

	passwordMatches := validatePassword(user.HashedPassword, params.Password)
	if !passwordMatches {
		return nil, errs.B().Code(errs.InvalidArgument).Msg("wrong email or password").Err()
	}

	session := Session{UserID: user.ID}
	err = s.db.Create(&session).Error
	if err != nil {
		// TODO: melhorar erro
		return nil, err
	}

	token, err := tokengenerator.GenerateTokenFor(tokengenerator.SessionToken, map[string]string{"SessionID": strconv.Itoa(int(session.ID))})
	if err != nil {
		return nil, err
	}

	response := LoginResponse{
		SessionToken: token,
	}

	return &response, nil
}
