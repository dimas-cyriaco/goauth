package user

import (
	"testing"
	"time"

	"encore.app/developer_area/internal/utils"
	"encore.dev/et"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RequestEmailVerificationSuite struct {
	UserTestSuite
}

func (suite *UserTestSuite) TestPublishToTopic() {
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

func (suite *UserTestSuite) TestDoNotPublishIfAlreadyVerified() {
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
