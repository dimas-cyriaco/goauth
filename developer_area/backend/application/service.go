package application

import (
	"time"

	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Application struct {
	ID           int       `json:"id"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	Name         string    `json:"name"`
	OwnerID      int       `json:"owner_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

//encore:service
type Service struct {
	db *gorm.DB
}

func initService() (*Service, error) {
	return NewApplicationService(db)
}

func NewApplicationService(database *sqldb.Database) (*Service, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: database.Stdlib()}))
	if err != nil {
		return nil, err
	}

	return &Service{db: db}, nil
}

var db = sqldb.NewDatabase("application", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
