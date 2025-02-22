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
	// TODO:
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
