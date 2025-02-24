package user

import (
	"strconv"
	"testing"

	"encore.app/developer_area/internal/tokens"
	"encore.app/developer_area/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	UserTestSuite
}

func (suite *AuthTestSuite) TestAuth() {
	// Arrange

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
	// Arrange

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
	// Arrange

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

// func (suite *AuthTestSuite) TestAuthShouldFailWithInvalidSessionToken() {
// 	// Arrange
//
// 	utils.Must(suite.RegisterUser())
//
// 	response := suite.Login()
// 	sessionCookie := findCookieByName(response.Result().Cookies(), "session_token")
// 	payload, _ := tokens.GetPayloadForToken(tokens.SessionToken, sessionCookie.Value)
//
// 	// TODO: Tamper with the session token payload.
//
// 	authData := AuthData{
// 		SessionToken: sessionCookie,
// 		CSRFToken:    payload["CSRFToken"],
// 	}
//
// 	// Act
//
// 	_, _, err := HandleAuthentication(suite.db, &authData)
//
// 	// Assert
//
// 	assert.Error(suite.T(), err)
// 	assert.ErrorContains(suite.T(), err, "unauthenticated")
// }

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
