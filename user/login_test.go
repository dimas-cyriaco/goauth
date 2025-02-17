package user

import (
	"context"
	"testing"

	"encore.app/utils"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LoginTestSuite struct {
	suite.Suite
	ctx     context.Context
	service *Service
}

func (suite *LoginTestSuite) SetupTest() {
	ctx := context.Background()

	service := utils.Must(initService())

	suite.ctx = ctx
	suite.service = service
}

func (suite *LoginTestSuite) TestLogin() {
	// Act

	password := "foo"
	email := faker.Email()

	a := RegistrationParams{
		Email:                email,
		Password:             password,
		PasswordConfirmation: password,
	}
	utils.Must(suite.service.Registration(suite.ctx, &a))

	// Act

	err := suite.service.Login(suite.ctx, &LoginParams{
		Email:    email,
		Password: password,
	})

	// Assert

	assert.NoError(suite.T(), err)
}

func (suite *LoginTestSuite) TestShouldFailWithWrongPassword() {
	// Act

	password := faker.Password()
	email := faker.Email()

	a := RegistrationParams{
		Email:                email,
		Password:             password,
		PasswordConfirmation: password,
	}
	utils.Must(suite.service.Registration(suite.ctx, &a))

	// Act

	err := suite.service.Login(suite.ctx, &LoginParams{
		Email:    faker.Email(),
		Password: "wrong-password",
	})

	// Assert

	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, "wrong email or password")
}

func (suite *LoginTestSuite) TestShouldFailWithWrongEmail() {
	// Act

	password := faker.Password()
	email := faker.Email()

	a := RegistrationParams{
		Email:                email,
		Password:             password,
		PasswordConfirmation: password,
	}
	utils.Must(suite.service.Registration(suite.ctx, &a))

	// Act

	err := suite.service.Login(suite.ctx, &LoginParams{
		Email:    "wrong@email.com",
		Password: password,
	})

	// Assert

	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, "wrong email or password")
}

func (suite *LoginTestSuite) TestShouldCreateSession()      {}
func (suite *LoginTestSuite) TestShouldReturnSessionToken() {}
func (suite *LoginTestSuite) TestShouldReturnCRSFToken()    {}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}
