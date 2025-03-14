package account

import (
	"fmt"
	"testing"

	"encore.app/oauth_flows/backend/account/db"
	"encore.dev/et"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockMailer is a mock implementation of the utils.Mailer interface
type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) SendEmail(to string, subject string, body string, config *MailerConfig) error {
	args := m.Called(to, subject, body, config)
	return args.Error(0)
}

type SendVerificationEmailTestSuit struct {
	AccountTestSuite
}

func (suite *SendVerificationEmailTestSuit) TestSendVerificationEmail() {
	// Arrange

	et.SetCfg(mailConfig.SendEmails, true)

	userID := Must(suite.RegisterAccount())

	mockMailer := new(MockMailer)

	link := Must(generateEmailVerificationLinkForUser(&db.Account{ID: userID}))
	emailBody := fmt.Sprintf("Welcome to GOAuth. To verify your email click this link: %s", link)

	mockMailer.On("SendEmail", suite.email, "Welcome to GOAuth", emailBody, &MailerConfig{
		SendEmailsFrom: mailConfig.SendEmailsFrom(),
		SMTPHost:       mailConfig.SMTPHost(),
		SMTPPort:       mailConfig.SMTPPort(),
		SMTPUsername:   mailConfig.SMTPUsername(),
		SMTPPassword:   secrets.SMTPPassword,
	}).Return(nil)

	// Act

	result := SendVerificationEmail(suite.ctx, &EmailVerificationRequestedEvent{
		UserID: userID,
	}, mockMailer, suite.db)

	// Assert

	assert.Nil(suite.T(), result)
}

func TestSendVerificationEmailTestSuite(t *testing.T) {
	suite.Run(t, new(SendVerificationEmailTestSuit))
}
