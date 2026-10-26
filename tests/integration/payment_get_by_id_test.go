package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPaymentByID_Success_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAdminToken(t, userDeps)
	userID, err := getUserIDFromToken(token)
	require.NoError(t, err)

	paymentID, err := insertTestPayment(context.Background(), paymentDeps.db, userID, TestPayment{
		Name:      "claude",
		Amount:    15.99,
		Type:      "expense",
		Category:  "work",
		Date:      "2026-01-03",
		Paid:      false,
		Recurrent: false,
	})
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/v1/admin/payments/"+strconv.Itoa(paymentID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response dto.PaymentResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, int32(paymentID), response.ID)
	assert.Equal(t, "claude", response.Name)
	assert.Equal(t, 15.99, response.Amount)
	assert.Equal(t, entity.TransactionTypeExpense, response.Type)
	assert.Equal(t, entity.CategoryTypeWork, response.Category)
}

func TestGetPaymentByID_InvalidID_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAdminToken(t, userDeps)

	req := httptest.NewRequest("GET", "/v1/admin/payments/abb", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid payment id", response.Error)
}

func TestGetPaymentByID_Unauthorized_Integration(t *testing.T) {
	_, paymentDeps := setupTestDependencies(t)

	req := httptest.NewRequest("GET", "/v1/admin/payments/1", nil)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "no token provided", response.Error)
}

func TestGetPaymentByID_NotFound_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAdminToken(t, userDeps)

	req := httptest.NewRequest("GET", "/v1/admin/payments/99999", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "payment not found", response.Error)
}

func TestGetPaymentByID_AdminCanAccessOtherUsersPayment_Integration(t *testing.T) {
	userDeps, paymentDeps := setupTestDependencies(t)
	token := getAdminToken(t, userDeps)
	token2 := getAdminToken(t, userDeps)
	userID, err := getUserIDFromToken(token)
	require.NoError(t, err)

	paymentID, err := insertTestPayment(context.Background(), paymentDeps.db, userID, TestPayment{
		Name:      "spotify",
		Amount:    9.99,
		Type:      "subscription",
		Category:  "entertainment",
		Date:      "2026-01-03",
		Paid:      false,
		Recurrent: false,
	})
	require.NoError(t, err)

	paymentIDStr := strconv.Itoa(paymentID)
	req := httptest.NewRequest("GET", "/v1/admin/payments/"+paymentIDStr, nil)
	req.Header.Set("Authorization", "Bearer "+token2)

	w := httptest.NewRecorder()
	paymentDeps.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response dto.PaymentResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, int32(paymentID), response.ID)
	assert.Equal(t, "spotify", response.Name)
}
