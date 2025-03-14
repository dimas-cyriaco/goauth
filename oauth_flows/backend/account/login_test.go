package account

import (
	"net/http"
	"testing"

	"encore.app/oauth_flows/backend/internal/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LoginTestSuite struct {
	AccountTestSuite
}

func (suite *LoginTestSuite) TestLogin() {
	// Arrange

	suite.RegisterAccount()

	// Act

	response := suite.Login()

	// Assert

	assert.Equal(suite.T(), http.StatusOK, response.Code)
	assert.NotEmpty(suite.T(), response.Result().Cookies())
}

func (suite *LoginTestSuite) TestShouldFailWithWrongPassword() {
	// Arrange

	suite.RegisterAccount()

	// Act

	response := suite.LoginWith(suite.email, "wrong-password")

	// Assert

	assert.Equal(suite.T(), http.StatusUnauthorized, response.Code)
	assert.Contains(suite.T(), response.Body.String(), "wrong email or password")
}

func (suite *LoginTestSuite) TestShouldFailWithWrongEmail() {
	// Arrange

	suite.RegisterAccount()

	// Act

	response := suite.LoginWith("wrong@email.com", suite.password)

	// Assert

	assert.Equal(suite.T(), http.StatusUnauthorized, response.Code)
	assert.Contains(suite.T(), response.Body.String(), "wrong email or password")
}

func (suite *LoginTestSuite) TestShouldCreateSession() {
	// Arrange

	suite.RegisterAccount()
	// var countBefore int64
	// TODO:
	// suite.service.db.Model(&Session{}).Count(&countBefore)

	// Act

	response := suite.Login()

	// Assert

	assert.Equal(suite.T(), http.StatusOK, response.Code)

	// var countAfter int64
	// suite.service.db.Model(&Session{}).Count(&countAfter)
	//
	// assert.Equal(suite.T(), countBefore+1, countAfter)
}

func (suite *LoginTestSuite) TestShouldReturnSessionToken() {
	// Arrange

	suite.RegisterAccount()

	// Act

	response := suite.Login()

	// Assert

	sessionCookie := findCookieByName(response.Result().Cookies(), "session_token")
	assert.NotNil(suite.T(), sessionCookie)

	_, err := tokens.GetPayloadForToken(tokens.SessionToken, sessionCookie.Value)
	assert.NoError(suite.T(), err)

	// TODO:
	// var session db.Session
	// suite.service.db.Model(&Session{}).Last(&session)
	// assert.Equal(suite.T(), payload["SessionID"], strconv.Itoa(int(session.ID)))
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}

func findCookieByName(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}
