package account

import (
	"encore.app/oauth_flows/backend/account/db"
	"encore.dev/storage/sqldb"
)

//encore:service
type Service struct {
	Query *db.Queries
}

func initService() (*Service, error) {
	return NewAccountService(accountsDB)
}

func NewAccountService(database *sqldb.Database) (*Service, error) {
	pgxdb := sqldb.Driver(database)
	query := db.New(pgxdb)

	return &Service{Query: query}, nil
}

var accountsDB = sqldb.NewDatabase("account", sqldb.DatabaseConfig{
	Migrations: "./db/migrations",
})
