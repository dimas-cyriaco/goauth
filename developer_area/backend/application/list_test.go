package application

import (
	"strconv"
	"testing"

	user_service "encore.app/developer_area/backend/user"
	"encore.app/developer_area/internal/utils"
	"encore.dev/beta/auth"
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
	userParams := user_service.RegistrationParams{
		Email:                email,
		Password:             password,
		PasswordConfirmation: password,
	}
	user := utils.Must(user_service.Registration(suite.ctx, &userParams))

	et.OverrideAuthInfo(auth.UID(strconv.Itoa(user.ID)), &user_service.AuthData{})

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
	userParams := user_service.RegistrationParams{
		Email:                email,
		Password:             password,
		PasswordConfirmation: password,
	}
	user := utils.Must(user_service.Registration(suite.ctx, &userParams))

	et.OverrideAuthInfo(auth.UID(strconv.Itoa(user.ID)), &user_service.AuthData{})

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
	userParams := user_service.RegistrationParams{
		Email:                email,
		Password:             password,
		PasswordConfirmation: password,
	}
	user := utils.Must(user_service.Registration(suite.ctx, &userParams))

	et.OverrideAuthInfo(auth.UID(strconv.Itoa(user.ID)), &user_service.AuthData{})

	appName := faker.Name()
	suite.service.Create(suite.ctx, &ApplicationParams{Name: appName})

	// Act

	response, err := suite.service.List(suite.ctx, &ApplicationListParams{Page: 1, PerPage: 10})

	// Assert

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), response.Applications[0].Name, appName)
}

// TODO: Test owner
// TODO: Test min page
// TODO: Test max page
// TODO: Test min per_page
// TODO: Test max per_page

func TestAppListTestSuite(t *testing.T) {
	suite.Run(t, new(ALSuite))
}
