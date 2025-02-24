package tokens

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"encore.app/developer_area/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TokensTestSuit struct {
	suite.Suite
	ctx context.Context
}

func (suite *TokensTestSuit) SetupTest() {
	ctx := context.Background()

	suite.ctx = ctx
}

func (suite *TokensTestSuit) TestGenerateTokenFor() {
	// Arrange

	purpose := PasswordRecovery
	payload := map[string]string{"userId": "123"}

	// Act

	token, _ := GenerateTokenFor(purpose, payload)

	// Assert

	assert.NotNil(suite.T(), token)

	decodedToken, _ := base64.URLEncoding.DecodeString(token)
	tokenString := string(decodedToken)

	assert.Equal(suite.T(), "{\"pur\":\"1\",\"userId\":\"123\"}", tokenString)
}

func (suite *TokensTestSuit) TestVerifyToken() {
	// Arrange

	purpose := PasswordRecovery
	payload := map[string]string{"userId": "123"}

	token, _ := GenerateTokenFor(purpose, payload)

	// Act

	isValid, _ := VerifyToken(purpose, token)

	// Assert

	assert.True(suite.T(), isValid)
}

func (suite *TokensTestSuit) TestProtectsFromTampering() {
	// Arrange:
	// Creates an token and extracts it's digest. Then, creates a token with a
	// diferent payload, and extracts it's encoded payload.

	purpose := PasswordRecovery

	originalPayload := map[string]string{"userId": "123"}
	correctToken := utils.Must(GenerateTokenFor(purpose, originalPayload))
	correctTokenSegments := utils.Must(GetTokenSegments(correctToken))
	originalDigest := correctTokenSegments.Digest

	manipulatedPayload := map[string]string{"userId": "456"}
	manipulatedToken := utils.Must(GenerateTokenFor(purpose, manipulatedPayload))
	manipulatedTokenSegments := utils.Must(GetTokenSegments(manipulatedToken))
	manipulatedEncodedPayload := manipulatedTokenSegments.EncodedPayload

	// Act:
	// Verifys a "new" token created by using the manipulated encoded payload
	// with the original digest.

	token := fmt.Sprintf("%s%s%s", manipulatedEncodedPayload, SEPARATOR, originalDigest)
	isValid := utils.Must(VerifyToken(purpose, token))

	// Assert

	assert.False(suite.T(), isValid)
}

func (suite *TokensTestSuit) TestGetPayload() {
	// Arrange

	purpose := PasswordRecovery
	payload := map[string]string{"userId": "123"}

	token := utils.Must(GenerateTokenFor(purpose, payload))

	// Act

	returnedPayload, err := GetPayloadForToken(purpose, token)

	// Assert

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), payload, returnedPayload)
}

func (suite *TokensTestSuit) TestGetPayloadValidatesPurpose() {
	// Arrange

	purpose := PasswordRecovery
	payload := map[string]string{"userId": "123"}

	token := utils.Must(GenerateTokenFor(purpose, payload))

	// Act

	wrongPurpose := EmailVerification
	_, err := GetPayloadForToken(wrongPurpose, token)

	// Assert

	assert.Error(suite.T(), err)
}

func TestTokensTestSuit(t *testing.T) {
	suite.Run(t, new(TokensTestSuit))
}
