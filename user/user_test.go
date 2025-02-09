package user

import (
	"context"
	"testing"

	"encore.app/utils"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserTestSuit struct {
	suite.Suite
	ctx     context.Context
	service *Service
}

func (suite *UserTestSuit) SetupTest() {
	ctx := context.Background()

	// WARN: Don't work!
	// et.NewTestDatabase(ctx, "user")

	service := utils.Must(initService())

	suite.ctx = ctx
	suite.service = service
}

func (suite *UserTestSuit) TestCreateUser() {
	// Act

	utils.Must(suite.service.Create(suite.ctx, &CreateParams{
		Name:  faker.Name(),
		Email: faker.Email(),
	}))

	// Assert

	user, _ := suite.service.Get(suite.ctx, 1)
	assert.NotNil(suite.T(), user)
}

func (suite *UserTestSuit) TestCreateUserValidatesPresenceOfName() {
	// Arrange

	params := CreateParams{
		Email: faker.Email(),
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)
	assert.Equal(suite.T(), "Key: 'CreateParams.Name' Error:Field validation for 'Name' failed on the 'required' tag", validationError.Error())
}

func (suite *UserTestSuit) TestCreateUserValidatesPresenceOfEmail() {
	// Arrange

	params := CreateParams{
		Name: faker.Word(),
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)
	assert.Equal(suite.T(), "Key: 'CreateParams.Email' Error:Field validation for 'Email' failed on the 'required' tag", validationError.Error())
}

func (suite *UserTestSuit) TestCreateUserValidatesFormatOfEmail() {
	// Arrange

	params := CreateParams{
		Name:  faker.Word(),
		Email: faker.Word(),
	}

	// Act

	validationError := params.Validate()

	// Assert

	assert.NotNil(suite.T(), validationError)
	assert.Equal(suite.T(), "Key: 'CreateParams.Email' Error:Field validation for 'Email' failed on the 'email' tag", validationError.Error())
}

func (suite *UserTestSuit) TestGetUser() {
	// Arrange

	createdUser := utils.Must(suite.service.Create(suite.ctx, &CreateParams{
		Name: faker.Name(),
	}))

	// Act

	user := utils.Must(suite.service.Get(suite.ctx, createdUser.ID))

	// Assert

	assert.Equal(suite.T(), user.Name, createdUser.Name)
}

func TestCreateUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuit))
}
