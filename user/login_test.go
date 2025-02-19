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

type LoginTestSuite struct {
	suite.Suite
	ctx      context.Context
	service  *Service
	email    string
	password string
}

func (suite *LoginTestSuite) SetupTest() {
	ctx := context.Background()

	service := utils.Must(initService())

	suite.ctx = ctx
	suite.service = service

	suite.password = faker.Password()
	suite.email = faker.Email()
}

//
// func (suite *LoginTestSuite) TeardownTest() {
// 	suite.password = ""
// 	suite.email = ""
// }

func (suite *LoginTestSuite) TestLogin() {
	// Act

	suite.registerUser()

	// Act

	_, err := suite.service.Login(suite.ctx, &LoginParams{
		Email:    suite.email,
		Password: suite.password,
	})

	// Assert

	assert.NoError(suite.T(), err)
}

func (suite *LoginTestSuite) TestShouldFailWithWrongPassword() {
	// Act

	suite.registerUser()

	// Act

	_, err := suite.service.Login(suite.ctx, &LoginParams{
		Email:    suite.email,
		Password: "wrong-password",
	})

	// Assert

	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, "wrong email or password")
}

func (suite *LoginTestSuite) TestShouldFailWithWrongEmail() {
	// Act

	suite.registerUser()

	// Act

	_, err := suite.service.Login(suite.ctx, &LoginParams{
		Email:    "wrong@email.com",
		Password: suite.password,
	})

	// Assert

	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, "wrong email or password")
}

func (suite *LoginTestSuite) TestShouldCreateSession() {
	// Act

	suite.registerUser()

	var countBefore int64
	suite.service.db.Model(&Session{}).Count(&countBefore)

	// Act

	utils.Must(suite.service.Login(suite.ctx, &LoginParams{
		Email:    suite.email,
		Password: suite.password,
	}))

	// Assert

	var countAfter int64
	suite.service.db.Model(&Session{}).Count(&countAfter)

	assert.Equal(suite.T(), countBefore, countAfter-1)
}

func (suite *LoginTestSuite) TestShouldReturnSessionToken() {
	// Act

	suite.registerUser()

	// Act

	result := utils.Must(suite.service.Login(suite.ctx, &LoginParams{
		Email:    suite.email,
		Password: suite.password,
	}))

	// Assert

	assert.NotNil(suite.T(), result.SessionToken)

	payload, _ := tokengenerator.GetPayloadForToken(tokengenerator.SessionToken, result.SessionToken)

	var session Session
	suite.service.db.Model(&Session{}).Last(&session)

	assert.Equal(suite.T(), payload["SessionID"], strconv.Itoa(session.ID))
}

func (suite *LoginTestSuite) TestShouldReturnCRSFToken() {
	// Act

	suite.registerUser()

	// Act

	result := utils.Must(suite.service.Login(suite.ctx, &LoginParams{
		Email:    suite.email,
		Password: suite.password,
	}))

	// Assert

	assert.NotEmpty(suite.T(), result.CSRFToken)

	sessionPayload, _ := tokengenerator.GetPayloadForToken(tokengenerator.SessionToken, result.SessionToken)
	assert.Equal(suite.T(), result.CSRFToken, sessionPayload["CSRFToken"])
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}

func (suite *LoginTestSuite) registerUser() {
	password := suite.password
	email := suite.email

	a := RegistrationParams{
		Email:                email,
		Password:             password,
		PasswordConfirmation: password,
	}

	utils.Must(suite.service.Registration(suite.ctx, &a))
}
