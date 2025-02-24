package user

import (
	"strconv"
	"testing"

	"encore.app/internal/tokens"
	"encore.app/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	UserTestSuite
}

func (suite *AuthTestSuite) TestAuth() {
	// Act

	userID := utils.Must(suite.RegisterUser())

	response := suite.Login()
	sessionCookie := findCookieByName(response.Result().Cookies(), "session_token")
	payload, _ := tokens.GetPayloadForToken(tokens.SessionToken, sessionCookie.Value)

	authData := AuthData{
		SessionToken: sessionCookie,
		CSRFToken:    payload["CSRFToken"],
	}

	// Act

	uid, _, err := HandleAuthentication(suite.db, &authData)

	// Assert

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), strconv.Itoa(userID), string(uid))
}

func (suite *AuthTestSuite) TestAuthShouldFailWithoutSessionToken() {
	// Act

	utils.Must(suite.RegisterUser())

	response := suite.Login()
	sessionCookie := findCookieByName(response.Result().Cookies(), "session_token")
	payload, _ := tokens.GetPayloadForToken(tokens.SessionToken, sessionCookie.Value)

	authData := AuthData{
		SessionToken: nil,
		CSRFToken:    payload["CSRFToken"],
	}

	// Act

	_, _, err := HandleAuthentication(suite.db, &authData)

	// Assert

	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, "unauthenticated")
}

func (suite *AuthTestSuite) TestAuthShouldFailWithoutCSRFToken() {
	// Act

	utils.Must(suite.RegisterUser())

	response := suite.Login()
	sessionCookie := findCookieByName(response.Result().Cookies(), "session_token")

	authData := AuthData{
		SessionToken: sessionCookie,
		CSRFToken:    "",
	}

	// Act

	_, _, err := HandleAuthentication(suite.db, &authData)

	// Assert

	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, "unauthenticated")
}

func (suite *AuthTestSuite) TestAuthShouldFailWithInvalidSessionToken() {
	// Act

	utils.Must(suite.RegisterUser())

	response := suite.Login()
	sessionCookie := findCookieByName(response.Result().Cookies(), "session_token")
	payload, _ := tokens.GetPayloadForToken(tokens.SessionToken, sessionCookie.Value)

	// TODO: Tamper with the session token payload.

	authData := AuthData{
		SessionToken: sessionCookie,
		CSRFToken:    payload["CSRFToken"],
	}

	// Act

	_, _, err := HandleAuthentication(suite.db, &authData)

	// Assert

	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, "unauthenticated")
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
