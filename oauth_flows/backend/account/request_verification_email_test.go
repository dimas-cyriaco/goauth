package account

import (
	"testing"

	"encore.app/oauth_flows/backend/account/db"
	"encore.dev/et"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RequestEmailVerificationSuite struct {
	AccountTestSuite
}

func (suite *RequestEmailVerificationSuite) TestPublishToTopic() {
	// Arrange

	insertParams := db.InsertAccountParams{}
	faker.FakeData(insertParams)
	Must(suite.service.Query.InsertAccount(suite.ctx, insertParams))

	params := RequestVerificationEmailParams{
		Email: insertParams.Email,
	}

	messagesBefore := et.Topic(EmailVerificationRequested).PublishedMessages()
	messageCountBefore := len(messagesBefore)

	// Act

	Must(suite.service.RequestVerificationEmail(suite.ctx, &params), nil)

	// Assert

	messagesAfter := et.Topic(EmailVerificationRequested).PublishedMessages()
	assert.Equal(suite.T(), messageCountBefore+1, len(messagesAfter))
}

func (suite *RequestEmailVerificationSuite) TestDoNotPublishIfAlreadyVerified() {
	// Arrange

	insertParams := db.InsertAccountParams{}
	faker.FakeData(&insertParams)
	acc := Must(suite.service.Query.InsertAccount(suite.ctx, insertParams))
	Must(suite.service.Query.VerifyEmail(suite.ctx, acc.ID), nil)

	params := RequestVerificationEmailParams{
		Email: insertParams.Email,
	}

	messagesBefore := et.Topic(EmailVerificationRequested).PublishedMessages()
	messageCountBefore := len(messagesBefore)

	// Act

	Must(suite.service.RequestVerificationEmail(suite.ctx, &params), nil)

	// Assert

	messagesAfter := et.Topic(EmailVerificationRequested).PublishedMessages()
	assert.Equal(suite.T(), messageCountBefore, len(messagesAfter))
}

func TestRequestEmailVerificationTestSuite(t *testing.T) {
	suite.Run(t, new(RequestEmailVerificationSuite))
}
