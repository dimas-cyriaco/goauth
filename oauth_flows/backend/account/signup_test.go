package account

import (
	"testing"

	"encore.app/internal/validation"
	"encore.dev/beta/errs"
	"encore.dev/et"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AccSuite struct {
	AccountTestSuite
}

func (suite *AccSuite) TestCreatesUser() {
	// Act

	password := faker.Word()
	response := Must(suite.service.Signup(suite.ctx, &SignupParams{
		Email:                faker.Email(),
		Password:             password,
		PasswordConfirmation: password,
	}))

	// Assert

	user := Must(suite.service.Query.ByID(suite.ctx, response.ID))

	assert.NotNil(suite.T(), response)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), response.ID, user.ID)
}

func (suite *AccSuite) TestValidatesPresenceOfEmail() {
	// Arrange

	password := faker.UUIDDigit()
	params := SignupParams{
		// Without Email
		Password:             password,
		PasswordConfirmation: password,
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)

	expected := &validation.ValidationErrors{
		"email": {"Email is a required field"},
	}

	assert.Equal(suite.T(), expected, validationError.(*errs.Error).Details)
}

func (suite *AccSuite) TestValidatesFormatOfEmail() {
	// Arrange

	password := faker.UUIDDigit()
	params := SignupParams{
		Email:                faker.DomainName(), // Email with wrong format.
		Password:             password,
		PasswordConfirmation: password,
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)

	errors := validationError.(*errs.Error)

	expected := &validation.ValidationErrors{
		"email": {"Email must be a valid email address"},
	}

	assert.Equal(suite.T(), expected, errors.Details)
}

func (suite *AccSuite) TestValidatesPresenceOfPassword() {
	// Arrange

	params := SignupParams{
		Email: faker.Email(),
		// No Password
		PasswordConfirmation: faker.UUIDDigit(),
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)

	errors := validationError.(*errs.Error)

	expected := &validation.ValidationErrors{
		"password":              {"Password is a required field"},
		"password_confirmation": {"PasswordConfirmation must be equal to Password"},
	}

	assert.Equal(suite.T(), expected, errors.Details)
}

func (suite *AccSuite) TestValidatesPresenceOfPasswordConfirmation() {
	// Arrange

	params := SignupParams{
		Email:    faker.Email(),
		Password: faker.UUIDDigit(),
		// No PasswordConfirmation
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)

	errors := validationError.(*errs.Error)

	expected := &validation.ValidationErrors{
		"password_confirmation": {"PasswordConfirmation is a required field"},
	}

	assert.Equal(suite.T(), expected, errors.Details)
}

func (suite *AccSuite) TestValidatesPasswordConfirmationMatch() {
	// Arrange

	params := SignupParams{
		Email:                faker.Email(),
		Password:             faker.UUIDDigit(),
		PasswordConfirmation: "this-will-not-match-the-password",
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)

	errors := validationError.(*errs.Error)

	expected := &validation.ValidationErrors{
		"password_confirmation": {"PasswordConfirmation must be equal to Password"},
	}

	assert.Equal(suite.T(), expected, errors.Details)
}

func (suite *AccSuite) TestHashesPassword() {
	// Arrange

	params := SignupParams{}
	faker.FakeData(&params)

	// Act

	response := Must(suite.service.Signup(suite.ctx, &params))

	// Assert

	user := Must(suite.service.Query.ByID(suite.ctx, response.ID))

	assert.NotEmpty(suite.T(), user.HashedPassword)
	assert.NotEqual(suite.T(), params.Password, user.HashedPassword)
}

func (suite *AccSuite) TestRequiresEmailToBeUnique() {
	// Arrange

	password := faker.Word()
	params := SignupParams{
		Email:                faker.Email(),
		Password:             password,
		PasswordConfirmation: password,
	}
	Must(suite.service.Signup(suite.ctx, &params))

	// Act

	_, err := suite.service.Signup(suite.ctx, &params)

	// Assert

	assert.NotNil(suite.T(), err)

	errors := err.(*errs.Error)

	expected := &validation.ValidationErrors{"email": {"Email already taken"}}

	assert.Equal(suite.T(), expected, errors.Details)
}

func (suite *AccSuite) TestTrimsEmailAndPassword() {
	// Arrange

	password := faker.Word()
	email := faker.Email()
	params := SignupParams{
		Email:                "  " + email + "   ",
		Password:             "  " + password + "   ",
		PasswordConfirmation: password,
	}

	// Act

	Must(params.Validate(), nil)

	// Assert

	// user, _ := suite.service.Get(suite.ctx, 1)
	assert.Equal(suite.T(), params.Email, email)
	assert.Equal(suite.T(), params.Password, password)
	assert.Equal(suite.T(), params.PasswordConfirmation, password)
}

func (suite *AccSuite) TestPublishToTopic() {
	// Arrange

	password := faker.Word()
	params := SignupParams{
		Email:                faker.Email(),
		Password:             password,
		PasswordConfirmation: password,
	}

	// Act

	Must(suite.service.Signup(suite.ctx, &params))

	// Assert

	// Get all published messages on the EmailVerificationRequested topic from this test.
	msgs := et.Topic(EmailVerificationRequested).PublishedMessages()
	assert.Len(suite.T(), msgs, 1)
}

func TestRegistrationTestSuite(t *testing.T) {
	suite.Run(t, new(AccSuite))
}
