package application

import (
	"strconv"
	"testing"

	"encore.app/developer_area/backend/internal/utils"
	account_service "encore.app/oauth_flows/backend/account"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/et"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ALSuite struct {
	ApplicationTestSuite
}

func (suite *ALSuite) TestListApplications() {
	// Arrange

	email := faker.Email()
	password := faker.Password()
	userParams := account_service.SignupParams{
		Email:                email,
		Password:             password,
		PasswordConfirmation: password,
	}
	user := utils.Must(account_service.Signup(suite.ctx, &userParams))

	et.OverrideAuthInfo(auth.UID(strconv.Itoa(int(user.ID))), &account_service.AuthData{})

	appName := faker.Name()
	suite.service.Create(suite.ctx, &ApplicationParams{Name: appName})

	// Act

	response, err := suite.service.List(suite.ctx, &ApplicationListParams{Page: 1, PerPage: 10})

	// Assert

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), response.Applications[0].Name, appName)
}

func (suite *ALSuite) TestShouldReturnPaginated() {
	// Arrange

	email := faker.Email()
	password := faker.Password()
	userParams := account_service.SignupParams{
		Email:                email,
		Password:             password,
		PasswordConfirmation: password,
	}
	user := utils.Must(account_service.Signup(suite.ctx, &userParams))

	et.OverrideAuthInfo(auth.UID(strconv.Itoa(int(user.ID))), &account_service.AuthData{})

	firstAppName := faker.Name()
	secondAppName := faker.Name()
	suite.service.Create(suite.ctx, &ApplicationParams{Name: firstAppName})
	suite.service.Create(suite.ctx, &ApplicationParams{Name: secondAppName})

	// Act

	params := ApplicationListParams{
		Page:    2,
		PerPage: 1,
	}

	response, err := suite.service.List(suite.ctx, &params)

	// Assert

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), response.Applications[0].Name, secondAppName)
}

func (suite *ALSuite) TestSetDefaultParams() {
	// Act

	params := ApplicationListParams{}
	validationError := params.Validate()

	// Assert

	assert.Nil(suite.T(), validationError)
	assert.Equal(suite.T(), ApplicationListParams{
		Page:    1,
		PerPage: 10,
	}, params)
}

func (suite *ALSuite) TestFilterByOwnerID() {
	// Arrange

	email := faker.Email()
	password := faker.Password()
	userParams := account_service.SignupParams{
		Email:                email,
		Password:             password,
		PasswordConfirmation: password,
	}
	user := utils.Must(account_service.Signup(suite.ctx, &userParams))

	et.OverrideAuthInfo(auth.UID(strconv.Itoa(int(user.ID))), &account_service.AuthData{})

	appName := faker.Name()
	suite.service.Create(suite.ctx, &ApplicationParams{Name: appName})

	// Act

	response, err := suite.service.List(suite.ctx, &ApplicationListParams{Page: 1, PerPage: 10})

	// Assert

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), response.Applications[0].Name, appName)
}

func (suite *ALSuite) TestValidateMinPage() {
	// Assert

	params := ApplicationListParams{
		Page: -1,
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)

	errors := validationError.(*errs.Error)
	expected := &utils.ValidationErrors{
		"page": {"Page must be 1 or greater"},
	}
	assert.Equal(suite.T(), expected, errors.Details)
}

func (suite *ALSuite) TestValidateMinPerPage() {
	// Assert

	params := ApplicationListParams{
		PerPage: -1,
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)

	errors := validationError.(*errs.Error)
	expected := &utils.ValidationErrors{
		"per_page": {"PerPage must be 1 or greater"},
	}
	assert.Equal(suite.T(), expected, errors.Details)
}

func (suite *ALSuite) TestValidateMaxPerPage() {
	// Assert

	params := ApplicationListParams{
		PerPage: 101,
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)

	errors := validationError.(*errs.Error)
	expected := &utils.ValidationErrors{
		"per_page": {"PerPage must be 100 or less"},
	}
	assert.Equal(suite.T(), expected, errors.Details)
}

func TestAppListTestSuite(t *testing.T) {
	suite.Run(t, new(ALSuite))
}
