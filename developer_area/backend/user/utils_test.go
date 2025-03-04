package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"encore.app/developer_area/internal/utils"
	"encore.dev/et"
	"encore.dev/storage/sqldb"
	"github.com/charmbracelet/log"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	ctx      context.Context
	db       *sqldb.Database
	service  *Service
	email    string
	password string
}

func (suite *UserTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.db = utils.Must(et.NewTestDatabase(suite.ctx, "user"))
	suite.service = utils.Must(NewUserService(suite.db))

	suite.email = faker.Email()
	suite.password = faker.Password()
}

// RegisterUser creates a new user account using the suite's default credentials.
//
// It uses the default suite email and password to register a new user.
func (suite *UserTestSuite) RegisterUser() (int, error) {
	log.Infof("🪵 suite.email: %v\n", suite.email)
	log.Infof("🪵 suite.password: %v\n", suite.password)
	params := RegistrationParams{
		Email:                suite.email,
		Password:             suite.password,
		PasswordConfirmation: suite.password,
	}

	response, err := suite.service.Registration(suite.ctx, &params)
	return response.ID, err
}

// Login performs a login request using the suite's default credentials.
//
// To use different credentials, use the `(suite *UserTestSuite) LoginWith(email, password string)` method.
func (suite *UserTestSuite) Login() *httptest.ResponseRecorder {
	log.Infof("🪵 suite.email: %v\n", suite.email)
	log.Infof("🪵 suite.password: %v\n", suite.password)
	return suite.LoginWith(suite.email, suite.password)
}

// LoginWith performs a login request with the specified credentials.
// It creates and executes an HTTP POST request to the /login endpoint with the given email and password.
func (suite *UserTestSuite) LoginWith(email, password string) *httptest.ResponseRecorder {
	loginData := url.Values{
		"email":    []string{email},
		"password": []string{password},
	}

	loginForm := strings.NewReader(loginData.Encode())

	request := httptest.NewRequest(http.MethodPost, "/login", loginForm)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()

	suite.service.Login(response, request)

	return response
}

func (suite *UserTestSuite) findUserByID(userID int) (*User, error) {
	var user User
	err := suite.service.db.Model(&User{}).Where("id = $1", userID).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
