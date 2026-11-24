package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllPayments_Success_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)
	userID1, err := getUserIDFromToken(token)
	require.NoError(t, err)
	token2 := getAuthToken(t, userDeps)
	userID2, err := getUserIDFromToken(token2)
	require.NoError(t, err)

	_, _ = insertTestPayment(context.Background(), paymentDeps.db, userID1, TestPayment{
		Name:      "claude",
		Amount:    15.99,
		Type:      "subscription",
		Category:  "work",
		Date:      "2026-01-03",
		Paid:      false,
		Recurrent: false,
	})
	_, _ = insertTestPayment(context.Background(), paymentDeps.db, userID1, TestPayment{
		Name:      "mac",
		Amount:    1700,
		Type:      "expense",
		Category:  "electronics",
		Date:      "2026-03-21",
		Paid:      true,
		Recurrent: false,
	})
	_, _ = insertTestPayment(context.Background(), paymentDeps.db, userID2, TestPayment{
		Name:      "work",
		Amount:    2300,
		Type:      "income",
		Category:  "work",
		Date:      "2026-01-03",
		Paid:      false,
		Recurrent: false,
	})

	adminToken := getAdminToken(t, userDeps)

	req := httptest.NewRequest("GET", "/v1/admin/payments", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response []dto.PaymentResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 3)

	assert.Contains(t, "claude", response[0].Name)
	assert.Contains(t, "mac", response[1].Name)
	assert.Contains(t, "work", response[2].Name)
	assert.False(t, response[0].Paid)
	assert.True(t, response[1].Paid)
	assert.False(t, response[2].Paid)
}

func TestGetAllPayments_EmptyList_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	adminToken := getAdminToken(t, userDeps)

	req := httptest.NewRequest("GET", "/v1/admin/payments", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response []dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 0)
}

func TestGetAllPayments_Unauthorized_Integration(t *testing.T) {
	_, paymentDeps := setupTestDependencies(t)

	req := httptest.NewRequest("GET", "/v1/admin/payments", nil)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "no token provided", response.Error)
}

func TestGetAllPayments_Forbidden_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAuthToken(t, userDeps)

	req := httptest.NewRequest("GET", "/v1/admin/payments", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "forbidden", response.Error)
}
