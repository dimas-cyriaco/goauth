package user

import (
	"time"

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
	HashedPassword  string     `encore:"sensitive"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" faker:"-"`
}

type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//encore:service
type Service struct {
	db *gorm.DB
}

func initService() (*Service, error) {
	return NewUserService(db)
}

func NewUserService(database *sqldb.Database) (*Service, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: database.Stdlib()}))
	if err != nil {
		return nil, err
	}

	return &Service{db: db}, nil
}

var db = sqldb.NewDatabase("user", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
