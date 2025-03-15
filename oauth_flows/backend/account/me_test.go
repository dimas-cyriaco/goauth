package account

import (
	"strconv"
	"testing"

	"encore.dev/beta/auth"
	"encore.dev/et"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MeTestSuite struct {
	AccountTestSuite
}

func (suite *MeTestSuite) TestShouldReturnCurrentLoggedInUserData() {
	// Arrange

	userID := Must(suite.RegisterAccount())

	et.OverrideAuthInfo(auth.UID(strconv.Itoa(int(userID))), &AuthData{})

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
