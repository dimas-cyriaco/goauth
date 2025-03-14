package account

import (
	"testing"

	"encore.app/oauth_flows/backend/account/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VerifyEmailTestSuit struct {
	AccountTestSuite
}

func (suite *VerifyEmailTestSuit) TestVerifyEmail() {
	// Arrange

	accountID := Must(suite.RegisterAccount())

	account := Must(suite.service.Query.ByID(suite.ctx, accountID))

	assert.Zero(suite.T(), account.EmailVerifiedAt)

	token := Must(generateEmailVerificationTokenForAccount(&db.Account{
		ID: accountID,
	}))

	// Act

	params := VerifyEmailParams{Token: token}
	err := suite.service.VerifyEmail(suite.ctx, &params)

	// Assert

	assert.NoError(suite.T(), err)

	account = Must(suite.service.Query.ByID(suite.ctx, accountID))

	assert.NotNil(suite.T(), account.EmailVerifiedAt)
}

func TestVerifyEmailTestSuite(t *testing.T) {
	suite.Run(t, new(VerifyEmailTestSuit))
}
