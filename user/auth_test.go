package user

import (
	"strconv"
	"testing"

	tokengenerator "encore.app/internal/token_generator"
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
	payload, _ := tokengenerator.GetPayloadForToken(tokengenerator.SessionToken, sessionCookie.Value)

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

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
