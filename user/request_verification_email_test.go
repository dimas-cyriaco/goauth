package user

import (
	"context"
	"testing"
	"time"

	"encore.app/utils"
	"encore.dev/et"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RequestEmailVerificationSuite struct {
	suite.Suite
	ctx     context.Context
	service *Service
}

func (suite *RequestEmailVerificationSuite) SetupTest() {
	ctx := context.Background()

	service := utils.Must(initService())

	suite.ctx = ctx
	suite.service = service
}

func (suite *UserTestSuit) TestPublishToTopic() {
	// Arrange

	user := &User{}
	faker.FakeData(user)
	utils.Must(suite.service.db.Create(user), nil)

	params := RequestVerificationEmailParams{
		Email: user.Email,
	}

	messagesBefore := et.Topic(EmailVerificationRequested).PublishedMessages()
	messageCountBefore := len(messagesBefore)

	// Act

	utils.Must(suite.service.RequestVerificationEmail(suite.ctx, &params), nil)

	// Assert

	messagesAfter := et.Topic(EmailVerificationRequested).PublishedMessages()
	assert.Equal(suite.T(), messageCountBefore+1, len(messagesAfter))
}

func (suite *UserTestSuit) TestDoNotPublishIfAlreadyVerified() {
	// Arrange

	user := User{
		EmailVerifiedAt: &time.Time{},
	}
	faker.FakeData(&user)
	utils.Must(suite.service.db.Create(&user), nil)

	params := RequestVerificationEmailParams{
		Email: user.Email,
	}

	messagesBefore := et.Topic(EmailVerificationRequested).PublishedMessages()
	messageCountBefore := len(messagesBefore)

	// Act

	utils.Must(suite.service.RequestVerificationEmail(suite.ctx, &params), nil)

	// Assert

	messagesAfter := et.Topic(EmailVerificationRequested).PublishedMessages()
	assert.Equal(suite.T(), messageCountBefore, len(messagesAfter))
}

func TestRequestEmailVerificationTestSuite(t *testing.T) {
	suite.Run(t, new(RequestEmailVerificationSuite))
}
