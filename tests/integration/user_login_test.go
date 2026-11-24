package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/rober0xf/notifier/internal/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin_Success_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)

	payload := `{
		"username": "rober",
		"email": "rober0xf@gmail.com",
		"password": "password1#!"
	}`

	registerReq := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(payload))
	registerReq.Header.Set("Content-Type", "application/json")

	registerW := httptest.NewRecorder()
	userDeps.router.ServeHTTP(registerW, registerReq)
	require.Equal(t, http.StatusCreated, registerW.Code)

	var created dto.UserCreatedResponse
	require.NoError(t, json.Unmarshal(registerW.Body.Bytes(), &created))
	userID := created.User.ID
	rawToken := insertVerificationToken(t, userDeps.db, userID)

	verifyReq := httptest.NewRequest("GET", "/v1/users/email_verification/"+rawToken, nil)
	verifyW := httptest.NewRecorder()
	userDeps.router.ServeHTTP(verifyW, verifyReq)
	require.Equal(t, http.StatusOK, verifyW.Code)

	loginBody := `{
		"email": "rober0xf@gmail.com",
		"password": "password1#!"
	}`
	loginReq := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")

	loginW := httptest.NewRecorder()
	userDeps.router.ServeHTTP(loginW, loginReq)
	require.Equal(t, http.StatusOK, loginW.Code)

	var response dto.LoginPayload
	err := json.Unmarshal(loginW.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "rober0xf@gmail.com", response.Email)

	cookies := loginW.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == session.SessionCookieName {
			sessionCookie = c
			break
		}
	}
	require.NotNil(t, sessionCookie, "session cookie should be set")
	assert.NotEmpty(t, sessionCookie.Value)
	assert.True(t, sessionCookie.HttpOnly)
}

func TestLogin_MissingFields_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)

	payload := `{
		"username": "rober0xf",
		"password": "password1#!"
	}`

	req := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "email and password are required", response.Error)
}

func TestLogin_InvalidCredentials_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)

	_, err := insertTestUser(context.Background(), userDeps.db, "test@gmail.com", "testuser", "password123#!")
	require.NoError(t, err)

	login := `{
		"email": "test@gmail.com",
		"password": "nopassword13#!"
	}`
	req := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(login))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	userDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid credentials", response.Error)
}

func TestLogin_EmailNotVerified_Integration(t *testing.T) {
	userDeps, _ := setupTestDependencies(t)

	payload := `{
		"username": "rober0xf",
		"email": "rober0xf@gmail.com",
		"password": "password1#!"
	}`

	registerReq := httptest.NewRequest("POST", "/v1/users/register", strings.NewReader(payload))
	registerReq.Header.Set("Content-Type", "application/json")

	registerW := httptest.NewRecorder()
	userDeps.router.ServeHTTP(registerW, registerReq)
	require.Equal(t, http.StatusCreated, registerW.Code)

	loginBody := `{
		"email": "rober0xf@gmail.com",
		"password": "password1#!"
	}`
	loginReq := httptest.NewRequest("POST", "/v1/users/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")

	loginW := httptest.NewRecorder()
	userDeps.router.ServeHTTP(loginW, loginReq)
	require.Equal(t, http.StatusForbidden, loginW.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(loginW.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "email not verified", response.Error)
}
