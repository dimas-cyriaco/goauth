package user

import (
	"context"
)

//encore:api public method=GET path=/users/:id
func (s *Service) Get(ctx context.Context, id int) (*User, error) {
	var user User
	if err := s.db.Where("id = $1", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
