package user

import (
	"context"
	"testing"

	"encore.app/utils"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VerifyEmailTestSuit struct {
	suite.Suite
	ctx     context.Context
	service *Service
}

func (suite *VerifyEmailTestSuit) SetupTest() {
	ctx := context.Background()

	service := utils.Must(initService())

	suite.ctx = ctx
	suite.service = service
}

func (suite *VerifyEmailTestSuit) TestVerifyEmail() {
	// Arrange

	a := RegistrationParams{}
	faker.FakeData(&a)
	response := utils.Must(suite.service.Registration(suite.ctx, &a))

	fromDB := utils.Must(suite.service.Get(suite.ctx, response.ID))

	assert.Nil(suite.T(), fromDB.EmailVerifiedAt)

	token := utils.Must(generateEmailVerificationTokenForUser(&User{
		ID: response.ID,
	}))

	// Act

	params := VerifyEmailParams{
		Token: token,
	}
	err := suite.service.VerifyEmail(suite.ctx, &params)

	// Assert

	assert.NoError(suite.T(), err)

	fromDB = utils.Must(suite.service.Get(suite.ctx, response.ID))

	assert.NotNil(suite.T(), fromDB.EmailVerifiedAt)
}

func TestVerifyEmailTestSuite(t *testing.T) {
	suite.Run(t, new(VerifyEmailTestSuit))
}
