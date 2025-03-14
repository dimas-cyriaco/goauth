package user

import (
	"context"
	"net/http"
	"strconv"

	"encore.app/developer_area/backend/internal/tokens"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AuthData struct {
	SessionToken *http.Cookie `cookie:"session_token"`
	CSRFToken    string       `header:"X-CSRF-Token"`
}

//encore:authhandler
func AuthHandler(ctx context.Context, data *AuthData) (auth.UID, *AuthData, error) {
	return HandleAuthentication(db, data)
}

func HandleAuthentication(database *sqldb.Database, data *AuthData) (auth.UID, *AuthData, error) {
	if data.SessionToken == nil || data.CSRFToken == "" {
		return auth.UID(""), data, &errs.Error{Code: errs.Unauthenticated, Message: "Invalid SessionToken or CSRFToken"}
	}

	sessionPayload, err := tokens.GetPayloadForToken(tokens.SessionToken, data.SessionToken.Value)
	if err != nil {
		return auth.UID(""), data, err
	}

	tokenMatches := data.CSRFToken == sessionPayload["CSRFToken"]
	if !tokenMatches {
		return auth.UID(""), data, &errs.Error{Code: errs.Unauthenticated, Message: "Invalid CSRFToken"}
	}

	gorm, err := gorm.Open(postgres.New(postgres.Config{Conn: database.Stdlib()}))
	if err != nil {
		return auth.UID(""), data, err
	}

	sessionID := sessionPayload["SessionID"]

	var session Session
	err = gorm.Where("id = $1", sessionID).First(&session).Error
	if err != nil {
		return auth.UID(strconv.Itoa(session.ID)), data, err
	}

	return auth.UID(strconv.Itoa(session.UserID)), data, nil
}
