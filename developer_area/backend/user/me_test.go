package user

import (
	"strconv"
	"testing"

	"encore.app/developer_area/internal/utils"
	"encore.dev/beta/auth"
	"encore.dev/et"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MeTestSuite struct {
	UserTestSuite
}

func (suite *MeTestSuite) TestShouldReturnCurrentLoggedInUserData() {
	// Arrange

	userID := utils.Must(suite.RegisterUser())

	et.OverrideAuthInfo(auth.UID(strconv.Itoa(userID)), &AuthData{})

	// Act

	me, err := suite.service.Me(suite.ctx)

	// Assert

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), me)
	assert.Equal(suite.T(), userID, me.ID)
}

func TestMeTestSuite(t *testing.T) {
	suite.Run(t, new(MeTestSuite))
}
