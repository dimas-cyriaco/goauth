package user

import (
	"context"
	"testing"

	"encore.app/utils"
	"encore.dev/et"
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
	et.NewTestDatabase(ctx, "TestCreateUser")

	service := utils.Must(initService())

	suite.ctx = ctx
	suite.service = service
}

func (suite *UserTestSuit) TestCreateUser() {
	// Act

	utils.Must(suite.service.Create(suite.ctx, &CreateParams{
		Name: "User Name",
	}))

	// Assert

	user, _ := suite.service.Get(suite.ctx, 1)
	assert.NotNil(suite.T(), user)
}

func (suite *UserTestSuit) TestGetUser() {
	// Arrange

	createdUser := utils.Must(suite.service.Create(suite.ctx, &CreateParams{
		Name: "User Name",
	}))

	// Act

	user := utils.Must(suite.service.Get(suite.ctx, createdUser.ID))

	// Assert

	assert.Equal(suite.T(), user.Name, createdUser.Name)
}

func TestCreateUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuit))
}
