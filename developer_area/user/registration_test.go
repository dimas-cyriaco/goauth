package user

import (
	"testing"

	"encore.app/developer_area/utils"
	"encore.dev/et"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RegistrationTestSuite struct {
	UserTestSuite
}

func (suite *RegistrationTestSuite) TestRegistration() {
	// Act

	password := faker.Word()
	response := utils.Must(suite.service.Registration(suite.ctx, &RegistrationParams{
		Email:                faker.Email(),
		Password:             password,
		PasswordConfirmation: password,
	}))

	// Assert

	user := utils.Must(suite.findUserByID(response.ID))

	assert.NotNil(suite.T(), response)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), response.ID, user.ID)
}

func (suite *RegistrationTestSuite) TestRegistrationValidatesPresenceOfEmail() {
	// Arrange

	password := faker.UUIDDigit()
	params := RegistrationParams{
		// Without Email
		Password:             password,
		PasswordConfirmation: password,
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)
	assert.Equal(suite.T(), "Key: 'RegistrationParams.Email' Error:Field validation for 'Email' failed on the 'required' tag", validationError.Error())
}

func (suite *RegistrationTestSuite) TestRegistrationValidatesFormatOfEmail() {
	// Arrange

	password := faker.UUIDDigit()
	params := RegistrationParams{
		Email:                faker.DomainName(), // Email with wrong format.
		Password:             password,
		PasswordConfirmation: password,
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)
	expectedError := "Key: 'RegistrationParams.Email' Error:Field validation for 'Email' failed on the 'email' tag"
	assert.Equal(suite.T(), expectedError, validationError.Error())
}

func (suite *RegistrationTestSuite) TestRegistrationValidatesPresenceOfPassword() {
	// Arrange

	params := RegistrationParams{
		Email: faker.Email(),
		// No Password
		PasswordConfirmation: faker.UUIDDigit(),
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)
	expectedError := "Key: 'RegistrationParams.Password' Error:Field validation for 'Password' failed on the 'required' tag"
	assert.Contains(suite.T(), validationError.Error(), expectedError)
}

func (suite *RegistrationTestSuite) TestRegistrationValidatesPresenceOfPasswordConfirmation() {
	// Arrange

	params := RegistrationParams{
		Email:    faker.Email(),
		Password: faker.UUIDDigit(),
		// No PasswordConfirmation
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)
	expectedError := "Key: 'RegistrationParams.PasswordConfirmation' Error:Field validation for 'PasswordConfirmation' failed on the 'required' tag"
	assert.Contains(suite.T(), validationError.Error(), expectedError)
}

func (suite *RegistrationTestSuite) TestRegistrationValidatesPasswordConfirmationMatch() {
	// Arrange

	params := RegistrationParams{
		Email:                faker.Email(),
		Password:             faker.UUIDDigit(),
		PasswordConfirmation: "this-will-not-match-the-password",
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)
	expectedError := "Key: 'RegistrationParams.PasswordConfirmation' Error:Field validation for 'PasswordConfirmation' failed on the 'eqcsfield' tag"
	assert.Contains(suite.T(), validationError.Error(), expectedError)
}

func (suite *RegistrationTestSuite) TestRegistrationHashesPassword() {
	// Arrange

	params := RegistrationParams{}
	faker.FakeData(&params)

	// Act

	response := utils.Must(suite.service.Registration(suite.ctx, &params))

	// Assert

	user := utils.Must(suite.findUserByID(response.ID))

	assert.NotEmpty(suite.T(), user.HashedPassword)
	assert.NotEqual(suite.T(), params.Password, user.HashedPassword)
}

func (suite *RegistrationTestSuite) TestRegistrationRequiresEmailToBeUnique() {
	// Arrange

	password := faker.Word()
	params := RegistrationParams{
		Email:                faker.Email(),
		Password:             password,
		PasswordConfirmation: password,
	}
	utils.Must(suite.service.Registration(suite.ctx, &params))

	// Act

	_, err := suite.service.Registration(suite.ctx, &params)

	// Assert

	assert.NotNil(suite.T(), err)
	assert.ErrorContains(suite.T(), err, "Invalid Argument")
}

func (suite *RegistrationTestSuite) TestRegistrationTrimsEmailAndPassword() {
	// Arrange

	password := faker.Word()
	email := faker.Email()
	params := RegistrationParams{
		Email:                "  " + email + "   ",
		Password:             "  " + password + "   ",
		PasswordConfirmation: password,
	}

	// Act

	utils.Must(params.Validate(), nil)

	// Assert

	// user, _ := suite.service.Get(suite.ctx, 1)
	assert.Equal(suite.T(), params.Email, email)
	assert.Equal(suite.T(), params.Password, password)
	assert.Equal(suite.T(), params.PasswordConfirmation, password)
}

func (suite *RegistrationTestSuite) TestRegistrationPublishToTopic() {
	// Arrange

	password := faker.Word()
	params := RegistrationParams{
		Email:                faker.Email(),
		Password:             password,
		PasswordConfirmation: password,
	}

	// Act

	utils.Must(suite.service.Registration(suite.ctx, &params))

	// Assert

	// Get all published messages on the EmailVerificationRequested topic from this test.
	msgs := et.Topic(EmailVerificationRequested).PublishedMessages()
	assert.Len(suite.T(), msgs, 1)
}

func TestRegistrationTestSuite(t *testing.T) {
	suite.Run(t, new(RegistrationTestSuite))
}
