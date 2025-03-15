package account

import (
	"context"
	"net/http"
	"strconv"

	"encore.app/oauth_flows/backend/account/db"
	"encore.app/oauth_flows/backend/internal/tokens"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
)

type AuthData struct {
	SessionToken *http.Cookie `cookie:"session_token"`
	CSRFToken    string       `header:"X-CSRF-Token"`
}

//encore:authhandler
func AuthHandler(ctx context.Context, data *AuthData) (auth.UID, *AuthData, error) {
	return HandleAuthentication(accountsDB, data)
}

func HandleAuthentication(database *sqldb.Database, data *AuthData) (auth.UID, *AuthData, error) {
	ctx := context.Background()

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

	sessionID, err := strconv.ParseInt(sessionPayload["SessionID"], 10, 64)
	if err != nil {
		rlog.Error("Error converting sessionID to int64.", "err", err)
		return auth.UID(""), data, &errs.Error{Code: errs.Unauthenticated, Message: "Invalid sessionID"}
	}

	pgxdb := sqldb.Driver(database)
	query := db.New(pgxdb)

	session, err := query.FindSessionByID(ctx, sessionID)
	if err != nil {
		return auth.UID(strconv.Itoa(int(session.ID))), data, err
	}

	return auth.UID(strconv.Itoa(int(session.AccountID))), data, nil
}
