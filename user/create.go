package user

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type CreateParams struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

//encore:api public method=POST path=/users
func (s *Service) Create(ctx context.Context, params *CreateParams) (*User, error) {
	user := &User{Name: params.Name}
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (req *CreateParams) Validate() error {
	validate := validator.New()
	return validate.Struct(req)
}
