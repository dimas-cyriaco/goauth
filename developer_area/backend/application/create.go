package application

import (
	"context"
	"errors"
	"strconv"

	"encore.app/developer_area/backend/internal/utils"
	"encore.app/internal/string_utils"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"github.com/jackc/pgx/v5/pgconn"
)

type ApplicationParams struct {
	Name string `mod:"trim" validate:"required" json:"name"`
}

type ApplicationResponse struct {
	ID           int    `json:"id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Name         string `json:"name"`
}

func (params *ApplicationParams) Validate() error {
	err := utils.ValidateTransform(context.Background(), params)
	return err
}

//encore:api auth method=POST path=/applications
func (s *Service) Create(ctx context.Context, params *ApplicationParams) (*ApplicationResponse, error) {
	userID, hasUserId := auth.UserID()
	if !hasUserId {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "unauthenticated"}
	}

	ownerID, err := strconv.Atoi(string(userID))
	if err != nil {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "unauthenticated"}
	}

	clientID, err := string_utils.GenerateSecureRandomString(32)
	if err != nil {
		return nil, &errs.Error{Code: errs.Unknown, Message: "Error generating secure clientID"}
	}

	clientSecret, err := string_utils.GenerateSecureRandomString(64)
	if err != nil {
		return nil, &errs.Error{Code: errs.Unknown, Message: "Error generating secure clientSecret"}
	}

	application := &Application{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Name:         params.Name,
		OwnerID:      ownerID,
	}

	err = s.db.Create(application).Error
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			details := &utils.ValidationErrors{
				"email": {"Application name already taken"},
			}
			return nil, &errs.Error{Code: errs.InvalidArgument, Message: "", Details: details}
		}

		rlog.Error("Database error creating new User.", "err", err)
		return nil, &errs.Error{Code: errs.Unknown, Message: "Unknown error"}
	}

	return &ApplicationResponse{
		ID:           application.ID,
		ClientID:     application.ClientID,
		ClientSecret: application.ClientSecret,
		Name:         application.Name,
	}, nil
}
