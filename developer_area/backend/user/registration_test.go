package user

import (
	"testing"

	"encore.app/developer_area/backend/internal/utils"
	"encore.dev/beta/errs"
	"encore.dev/et"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RSuite struct {
	UserTestSuite
}

func (suite *RSuite) TestCreatesUser() {
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

func (suite *RSuite) TestValidatesPresenceOfEmail() {
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

	expected := &utils.ValidationErrors{
		"email": {"Email is a required field"},
	}

	assert.Equal(suite.T(), expected, validationError.(*errs.Error).Details)
}

func (suite *RSuite) TestValidatesFormatOfEmail() {
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

	errors := validationError.(*errs.Error)

	expected := &utils.ValidationErrors{
		"email": {"Email must be a valid email address"},
	}

	assert.Equal(suite.T(), expected, errors.Details)
}

func (suite *RSuite) TestValidatesPresenceOfPassword() {
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

	errors := validationError.(*errs.Error)

	expected := &utils.ValidationErrors{
		"password":              {"Password is a required field"},
		"password_confirmation": {"PasswordConfirmation must be equal to Password"},
	}

	assert.Equal(suite.T(), expected, errors.Details)
}

func (suite *RSuite) TestValidatesPresenceOfPasswordConfirmation() {
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

	errors := validationError.(*errs.Error)

	expected := &utils.ValidationErrors{
		"password_confirmation": {"PasswordConfirmation is a required field"},
	}

	assert.Equal(suite.T(), expected, errors.Details)
}

func (suite *RSuite) TestValidatesPasswordConfirmationMatch() {
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

	errors := validationError.(*errs.Error)

	expected := &utils.ValidationErrors{
		"password_confirmation": {"PasswordConfirmation must be equal to Password"},
	}

	assert.Equal(suite.T(), expected, errors.Details)
}

func (suite *RSuite) TestHashesPassword() {
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

func (suite *RSuite) TestRequiresEmailToBeUnique() {
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

	errors := err.(*errs.Error)

	expected := &utils.ValidationErrors{"email": {"Email already taken"}}

	assert.Equal(suite.T(), expected, errors.Details)
}

func (suite *RSuite) TestTrimsEmailAndPassword() {
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

func (suite *RSuite) TestPublishToTopic() {
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
	suite.Run(t, new(RSuite))
}
