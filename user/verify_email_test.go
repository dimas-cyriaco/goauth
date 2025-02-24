package user

import (
	"testing"

	"encore.app/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VerifyEmailTestSuit struct {
	UserTestSuite
}

func (suite *VerifyEmailTestSuit) TestVerifyEmail() {
	// Arrange

	userID := utils.Must(suite.RegisterUser())

	fromDB := utils.Must(suite.service.Get(suite.ctx, userID))

	assert.Nil(suite.T(), fromDB.EmailVerifiedAt)

	token := utils.Must(generateEmailVerificationTokenForUser(&User{
		ID: userID,
	}))

	// Act

	params := VerifyEmailParams{Token: token}
	err := suite.service.VerifyEmail(suite.ctx, &params)

	// Assert

	assert.NoError(suite.T(), err)

	fromDB = utils.Must(suite.service.Get(suite.ctx, userID))

	assert.NotNil(suite.T(), fromDB.EmailVerifiedAt)
}

func TestVerifyEmailTestSuite(t *testing.T) {
	suite.Run(t, new(VerifyEmailTestSuit))
}
