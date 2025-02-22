package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	ctx      context.Context
	service  *Service
	email    string
	password string
}

func (suite *UserTestSuite) RegisterUser() (int, error) {
	a := RegistrationParams{
		Email:                suite.email,
		Password:             suite.password,
		PasswordConfirmation: suite.password,
	}

	response, err := suite.service.Registration(suite.ctx, &a)
	return response.ID, err
}

func (suite *UserTestSuite) Login() *httptest.ResponseRecorder {
	return suite.LoginWith(suite.email, suite.password)
}

func (suite *UserTestSuite) LoginWith(email, password string) *httptest.ResponseRecorder {
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}

	body, _ := json.Marshal(loginData)

	request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	suite.service.Login(response, request)

	return response
}
