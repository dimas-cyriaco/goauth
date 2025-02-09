package user

import (
	"context"
)

type CreateParams struct {
	Name string
}

//encore:api public method=POST path=/users
func (s *Service) Create(ctx context.Context, params *CreateParams) (*User, error) {
	user := &User{Name: params.Name}
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
