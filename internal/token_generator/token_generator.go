package tokengenerator

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"encore.dev/rlog"
)

var secrets struct {
	SecretMasterKey string // Secret used for cryptography.
}

type Purpose int

const (
	EmailVerification Purpose = 0
	PasswordRecovery  Purpose = 1
	SessionToken      Purpose = 2
	CSRFToken         Purpose = 3
)

const (
	SEPARATOR = "--"
)

type TokenSegments struct {
	EncodedPayload string
	Digest         string
}

// TODO: Add support for expiration date.

// GenerateTokenFor generates a signed token with the encoded payload for a
// given Purpose.
//
// The token is created with the format `{encoded-payload}--{digest}`.
//
// The encoded payload is not encrypted, only encoded with base64. It's
// contents is accessible by the user, but the *digest* garantees that the
// content of the payload cannot be tampered with.
func GenerateTokenFor(purpose Purpose, payload map[string]string) (string, error) {
	payload["pur"] = strconv.Itoa(int(purpose))

	marshalledPayload, err := json.Marshal(payload)
	if err != nil {
		rlog.Error("Error marshalling token payload.", "err", err)
		return "", err
	}

	encodedPayload := base64.URLEncoding.EncodeToString(marshalledPayload)

	digest := digestFor(encodedPayload)

	token := fmt.Sprintf("%s%s%s", encodedPayload, SEPARATOR, digest)

	return token, nil
}

func VerifyToken(purpose Purpose, token string) (bool, error) {
	segments, err := GetTokenSegments(token)
	if err != nil {
		return false, err
	}

	expectedDigest := digestFor(segments.EncodedPayload)

	comparisonResult := subtle.ConstantTimeCompare([]byte(segments.Digest), []byte(expectedDigest))
	isValid := comparisonResult == 1

	if !isValid {
		rlog.Error("Invalid token", "token", token)
		return false, err
	}

	payload, err := decodePayload(segments.EncodedPayload)
	if err != nil {
		return false, err
	}

	if payload["pur"] != strconv.Itoa(int(purpose)) {
		rlog.Error("Invalid token purpose", "token", token, "purpose", purpose)
		err := errors.New("INVALID PURPOSE")
		return false, err
	}

	return isValid, nil
}

func GetPayloadForToken(purpose Purpose, token string) (map[string]string, error) {
	isValid, err := VerifyToken(purpose, token)
	if err != nil {
		return map[string]string{}, err
	}

	if !isValid {
		err := errors.New("INVALID TOKEN")
		return map[string]string{}, err
	}

	segments, err := GetTokenSegments(token)
	if err != nil {
		return map[string]string{}, err
	}

	return decodePayload(segments.EncodedPayload)
}

func GetTokenSegments(token string) (*TokenSegments, error) {
	v := strings.Split(token, SEPARATOR)
	if len(v) != 2 {
		err := errors.New("INVALID TOKEN SEGMENT COUNT")
		rlog.Error("Error getting token segments", "token", token)
		return &TokenSegments{}, err
	}

	segments := TokenSegments{
		EncodedPayload: v[0],
		Digest:         v[1],
	}

	return &segments, nil
}

func decodePayload(encodedPayload string) (map[string]string, error) {
	decodedPayload, err := base64.URLEncoding.DecodeString(encodedPayload)
	if err != nil {
		rlog.Error("Error decoding token payload", "err", err, "payload", decodedPayload)
		return map[string]string{}, err
	}

	var unmarshalledPayload map[string]string
	err = json.Unmarshal(decodedPayload, &unmarshalledPayload)
	if err != nil {
		rlog.Error("Error unmarshalling token payload", "err", err, "payload", decodedPayload)
		return map[string]string{}, err
	}

	return unmarshalledPayload, nil
}

func digestFor(payload string) string {
	hash := hmac.New(sha256.New, []byte(secrets.SecretMasterKey))
	hash.Write([]byte(payload))

	return hex.EncodeToString(hash.Sum(nil))
}
