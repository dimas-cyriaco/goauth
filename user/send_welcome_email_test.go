package user

import (
	"context"
	"fmt"
	"testing"

	"encore.app/utils"
	"encore.dev/et"
	"github.com/go-faker/faker/v4"
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

type SendWelcomeEmailTestSuit struct {
	suite.Suite
	ctx                 context.Context
	registrationService *Service
}

func (suite *SendWelcomeEmailTestSuit) SetupTest() {
	ctx := context.Background()

	service := utils.Must(initService())

	suite.ctx = ctx
	suite.registrationService = service
}

func (suite *SendWelcomeEmailTestSuit) TestSendWelcomeEmail() {
	// Arrange

	et.SetCfg(mailConfig.SendEmails, true)

	a := RegistrationParams{}
	faker.FakeData(&a)
	response := utils.Must(suite.registrationService.Registration(suite.ctx, &a))

	mockMailer := new(MockMailer)

	link := utils.Must(generateEmailVerificationLinkForUser(&User{ID: response.ID}))
	emailBody := fmt.Sprintf("Welcome to GOAuth. To verify your email click this link: %s", link)

	mockMailer.On("SendEmail", a.Email, "Welcome to GOAuth", emailBody, &utils.MailerConfig{
		SendEmailsFrom: mailConfig.SendEmailsFrom(),
		SMTPHost:       mailConfig.SMTPHost(),
		SMTPPort:       mailConfig.SMTPPort(),
		SMTPUsername:   mailConfig.SMTPUsername(),
		SMTPPassword:   secrets.SMTPPassword,
	}).Return(nil)

	// Act

	result := SendWelcomeEmail(suite.ctx, &SignupEvent{
		UserID: response.ID,
	}, mockMailer)

	// Assert

	assert.Nil(suite.T(), result)
}

func TestSendWelcomeEmailTestSuite(t *testing.T) {
	suite.Run(t, new(SendWelcomeEmailTestSuit))
}
