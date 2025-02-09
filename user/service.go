package user

import (
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	// ID is a unique ID for the user.
	ID int `json:"id"`
	// Name is the user's Name.
	Name string `json:"name" encore:"sensitive"`
}

// This is a service struct, learn more: https://encore.dev/docs/go/primitives/service-structs
//
//encore:service
type Service struct {
	db *gorm.DB
}

// initService is automatically called by Encore when the service starts up.
func initService() (*Service, error) {
	bloft, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db.Stdlib(),
	}))
	if err != nil {
		return nil, err
	}
	return &Service{db: bloft}, nil
}

var db = sqldb.NewDatabase("user", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
