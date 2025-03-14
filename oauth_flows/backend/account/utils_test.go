package account

import (
	"context"

	"encore.dev/et"
	"encore.dev/storage/sqldb"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/suite"
)

type AccountTestSuite struct {
	suite.Suite
	ctx      context.Context
	db       *sqldb.Database
	service  *Service
	email    string
	password string
}

func GetAccountTestService(ctx context.Context) *Service {
	db := Must(et.NewTestDatabase(ctx, "account"))
	return Must(NewAccountService(db))
}

func (suite *AccountTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.db = Must(et.NewTestDatabase(suite.ctx, "account"))
	suite.service = Must(NewAccountService(suite.db))

	suite.email = faker.Email()
	suite.password = faker.Password()
}

// RegisterAccount creates a new Account account using the suite's default credentials.
//
// It uses the default suite email and password to register a new Account.
func (suite *AccountTestSuite) RegisterAccount() (int64, error) {
	params := SignupParams{
		Email:                suite.email,
		Password:             suite.password,
		PasswordConfirmation: suite.password,
	}

	response, err := suite.service.Signup(suite.ctx, &params)
	return response.ID, err
}

// Login performs a login request using the suite's default credentials.
//
// To use different credentials, use the `(suite *UserTestSuite) LoginWith(email, password string)` method.
// func (suite *UserTestSuite) Login() *httptest.ResponseRecorder {
// 	return suite.LoginWith(suite.email, suite.password)
// }

// LoginWith performs a login request with the specified credentials.
// It creates and executes an HTTP POST request to the /login endpoint with the given email and password.
// func (suite *UserTestSuite) LoginWith(email, password string) *httptest.ResponseRecorder {
// 	loginData := url.Values{
// 		"email":    []string{email},
// 		"password": []string{password},
// 	}
//
// 	loginForm := strings.NewReader(loginData.Encode())
//
// 	request := httptest.NewRequest(http.MethodPost, "/login", loginForm)
// 	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
// 	response := httptest.NewRecorder()
//
// 	suite.service.Login(response, request)
//
// 	return response
// }

// func (suite *UserTestSuite) findUserByID(userID int) (*User, error) {
// 	var user User
// 	err := suite.service.db.Model(&User{}).Where("id = $1", userID).First(&user).Error
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &user, nil
// }

func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}
