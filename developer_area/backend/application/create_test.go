package application

import (
	"strconv"
	"testing"

	"encore.app/developer_area/backend/internal/utils"
	account_service "encore.app/oauth_flows/backend/account"
	"encore.dev/beta/auth"
	"encore.dev/et"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ACSuite struct {
	ApplicationTestSuite
}

func (suite *ACSuite) TestCreatesApplication() {
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

	// Act

	applicationParams := ApplicationParams{Name: faker.Name()}
	application, err := suite.service.Create(suite.ctx, &applicationParams)

	// Assert

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), application)
	assert.NotEmpty(suite.T(), application.ClientID)
	assert.NotEmpty(suite.T(), application.ClientSecret)
	assert.Equal(suite.T(), applicationParams.Name, application.Name)
}

// Test Pagination
func (suite *ACSuite) TestShouldReturnPage() {
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

	applicationParams := ApplicationParams{Name: faker.Name()}
	application, err := suite.service.Create(suite.ctx, &applicationParams)

	// Assert

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), application)
	assert.Equal(suite.T(), applicationParams.Name, application.Name)
}

func TestAppCreateTestSuite(t *testing.T) {
	suite.Run(t, new(ACSuite))
}
