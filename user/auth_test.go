package user

import (
	"context"
	"strconv"
	"testing"

	tokengenerator "encore.app/internal/token_generator"
	"encore.app/utils"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	UserTestSuite
}

func (suite *AuthTestSuite) SetupTest() {
	ctx := context.Background()

	service := utils.Must(initService())

	suite.ctx = ctx
	suite.service = service

	suite.password = faker.Password()
	suite.email = faker.Email()
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
	ctx := context.Background()

	// Act

	uid, _, err := AuthHandler(ctx, &authData)

	// Assert

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), strconv.Itoa(userID), string(uid))
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
