package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings" // <-- IMPORT
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants" // <-- IMPORT
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

func TestBusinessConnectionAPI_Integration_Lifecycle(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	// --- USE THE HELPER (MUCH CLEANER) ---
	userA, tokenA := CreateTestUserAndLogin(t, router, "usera@conn.com", "ValidPass123!")
	userB, tokenB := CreateTestUserAndLogin(t, router, "userb@conn.com", "ValidPass123!")

	bizA := models.Business{Name: "Conn Biz A", OperatorUserID: userA.ID, BusinessType: "Other", BusinessCategory: "Mixed", BusinessPhase: "Growth"}
	testutil.TestDB.Create(&bizA)
	bizB := models.Business{Name: "Conn Biz B", OperatorUserID: userB.ID, BusinessType: "Other", BusinessCategory: "Mixed", BusinessPhase: "Growth"}
	testutil.TestDB.Create(&bizB)

	// --- USE CONSTANTS ---
	constConnectBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.ConnectBase
	constBizBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.BusinessBase

	var createdConn ports.BusinessConnectionResponse

	t.Run("Create Connection (User A)", func(t *testing.T) {
		createDTO := ports.CreateBusinessConnectionInput{
			InitiatingBusinessID: bizA.ID,
			ReceivingBusinessID:  bizB.ID,
			ConnectionType:       models.ConnectionTypePartnership,
			InitiatedByUserID:    userA.ID,
		}
		body, _ := json.Marshal(createDTO)

		req, _ := http.NewRequest(http.MethodPost, constConnectBase, bytes.NewBuffer(body)) // Use constant
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenA)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		json.Unmarshal(w.Body.Bytes(), &createdConn)
		assert.Equal(t, bizA.ID, createdConn.InitiatingBusiness.ID)
		assert.Equal(t, models.ConnectionStatusPending, createdConn.Status)
	})

	t.Run("Get Connections for Biz B (User B)", func(t *testing.T) {
		// Build path safely from constants
		bizBConnectPath := strings.Replace(constants.AppRoutes.BusinessConnects, ":id", fmt.Sprintf("%d", bizB.ID), 1)
		url := fmt.Sprintf("%s%s?status=pending", constBizBase, bizBConnectPath)

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+tokenB)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var connections ports.BusinessConnectionsResponse
		json.Unmarshal(w.Body.Bytes(), &connections)
		assert.Equal(t, 1, connections.Count)
	})

	t.Run("Accept Connection (User B)", func(t *testing.T) {
		// Build path safely from constants
		acceptPath := strings.Replace(constants.AppRoutes.ConnectAccept, ":id", fmt.Sprintf("%d", createdConn.ID), 1)
		url := constConnectBase + acceptPath

		req, _ := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Authorization", "Bearer "+tokenB)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var acceptedConn ports.BusinessConnectionResponse
		json.Unmarshal(w.Body.Bytes(), &acceptedConn)
		assert.Equal(t, models.ConnectionStatusActive, acceptedConn.Status)
	})
}
