package user

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"encore.app/developer_area/backend/internal/tokens"
)

//encore:api public raw method=POST path=/login
func (s *Service) Login(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	email := request.FormValue("email")
	password := request.FormValue("password")

	var user User
	err := s.db.
		Where("email = $1", email).
		First(&user).
		Error
	if err != nil {
		http.Error(response, "wrong email or password", http.StatusUnauthorized)
		return
	}

	passwordMatches := validatePassword(user.HashedPassword, password)
	if !passwordMatches {
		http.Error(response, "wrong email or password", http.StatusUnauthorized)
		return
	}

	session := Session{UserID: user.ID}
	err = s.db.Create(&session).Error
	if err != nil {
		http.Error(response, "wrong email or password", http.StatusUnauthorized)
		return
	}

	csrfToken, err := generateCSRFToken()
	if err != nil {
		http.Error(response, "Internal server error", http.StatusInternalServerError)
		return
	}

	sessionToken, err := tokens.GenerateTokenFor(tokens.SessionToken, map[string]string{
		"SessionID": strconv.Itoa(int(session.ID)),
		"CSRFToken": csrfToken,
	})
	if err != nil {
		http.Error(response, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set the session cookie
	http.SetCookie(response, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	// Set the csrf token cookie
	http.SetCookie(response, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	// Return success response
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(map[string]string{
		"status": "success",
	})
}

func generateCSRFToken() (string, error) {
	csrfBytes := make([]byte, 32)
	rand.Read(csrfBytes)
	csrfPayload := hex.EncodeToString(csrfBytes)

	return tokens.GenerateTokenFor(tokens.CSRFToken, map[string]string{
		"CSRFToken": csrfPayload,
	})
}
