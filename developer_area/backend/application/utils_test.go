package application

import (
	"context"

	"encore.app/developer_area/internal/utils"
	"encore.dev/et"
	"encore.dev/storage/sqldb"
	"github.com/stretchr/testify/suite"
)

type ApplicationTestSuite struct {
	suite.Suite
	ctx     context.Context
	db      *sqldb.Database
	service *Service
}

func (suite *ApplicationTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.db = utils.Must(et.NewTestDatabase(suite.ctx, "application"))
	suite.service = utils.Must(NewApplicationService(suite.db))
}
