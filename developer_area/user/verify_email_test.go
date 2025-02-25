package user

import (
	"testing"

	"encore.app/developer_area/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VerifyEmailTestSuit struct {
	UserTestSuite
}

func (suite *VerifyEmailTestSuit) TestVerifyEmail() {
	// Arrange

	userID := utils.Must(suite.RegisterUser())

	user := utils.Must(suite.findUserByID(userID))

	assert.Nil(suite.T(), user.EmailVerifiedAt)

	token := utils.Must(generateEmailVerificationTokenForUser(&User{
		ID: userID,
	}))

	// Act

	params := VerifyEmailParams{Token: token}
	err := suite.service.VerifyEmail(suite.ctx, &params)

	// Assert

	assert.NoError(suite.T(), err)

	user = utils.Must(suite.findUserByID(userID))

	assert.NotNil(suite.T(), user.EmailVerifiedAt)
}

func TestVerifyEmailTestSuite(t *testing.T) {
	suite.Run(t, new(VerifyEmailTestSuit))
}
