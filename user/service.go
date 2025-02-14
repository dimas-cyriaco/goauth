package user

import (
	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	// ID is a unique ID for the user.
	ID int `json:"id"`
	// Email is the user's Email.
	Email string `json:"email" encore:"sensitive"`
	// HashedPassword is the hashed version of the user's password.
	HashedPassword string `encore:"sensitive"`
}

//encore:service
type Service struct {
	db *gorm.DB
}

func initService() (*Service, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: db.Stdlib()}))
	if err != nil {
		return nil, err
	}

	return &Service{db: db}, nil
}

var db = sqldb.NewDatabase("user", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
