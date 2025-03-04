package user

import (
	"fmt"
	"testing"

	"encore.app/developer_area/internal/utils"
	"encore.dev/et"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockMailer is a mock implementation of the utils.Mailer interface
type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) SendEmail(to string, subject string, body string, config *utils.MailerConfig) error {
	args := m.Called(to, subject, body, config)
	return args.Error(0)
}

type SendVerificationEmailTestSuit struct {
	UserTestSuite
}

func (suite *SendVerificationEmailTestSuit) TestSendVerificationEmail() {
	// Arrange

	et.SetCfg(mailConfig.SendEmails, true)

	userID := utils.Must(suite.RegisterUser())

	mockMailer := new(MockMailer)

	link := utils.Must(generateEmailVerificationLinkForUser(&User{ID: userID}))
	emailBody := fmt.Sprintf("Welcome to GOAuth. To verify your email click this link: %s", link)

	mockMailer.On("SendEmail", suite.email, "Welcome to GOAuth", emailBody, &utils.MailerConfig{
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
