package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"

	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
)

func TestUserSessionService_Integration(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)

	userSessionService := services.NewUserSessionService(testutil.TestDB)

	user := models.User{FirstName: "Test", LoginEmail: "test@test.com", Active: true}
	testutil.TestDB.Create(&user)

	activeSession := models.UserSession{
		UserID:    user.ID,
		TokenHash: "active_token_hash",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	testutil.TestDB.Create(&activeSession)

	expiredSession := models.UserSession{
		UserID:    user.ID,
		TokenHash: "expired_token_hash",
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	}
	testutil.TestDB.Create(&expiredSession)

	revokedSession := models.UserSession{
		UserID:    user.ID,
		TokenHash: "revoked_token_hash",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		RevokedAt: &[]time.Time{time.Now()}[0],
	}
	testutil.TestDB.Create(&revokedSession)

	t.Run("Success - Get Session By ID", func(t *testing.T) {
		session, err := userSessionService.GetSessionByID(context.Background(), activeSession.ID)
		require.NoError(t, err)
		assert.Equal(t, activeSession.ID, session.ID)
		assert.Equal(t, user.ID, session.UserID)
		assert.Equal(t, "active_token_hash", session.TokenHash)
	})

	t.Run("Failure - Get Non-existent Session", func(t *testing.T) {
		_, err := userSessionService.GetSessionByID(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrSessionNotFound, err)
	})

	t.Run("Success - Get Active Sessions", func(t *testing.T) {
		sessions, err := userSessionService.GetActiveSessions(context.Background(), user.ID)
		require.NoError(t, err)
		assert.Len(t, sessions, 1)
		assert.Equal(t, activeSession.ID, sessions[0].ID)
	})

	t.Run("Success - Revoke Session", func(t *testing.T) {
		err := userSessionService.RevokeSession(context.Background(), activeSession.ID)
		require.NoError(t, err)

		session, err := userSessionService.GetSessionByID(context.Background(), activeSession.ID)
		require.NoError(t, err)
		assert.NotNil(t, session.RevokedAt)
	})

	t.Run("Failure - Revoke Non-existent Session", func(t *testing.T) {
		err := userSessionService.RevokeSession(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrSessionNotFound, err)
	})

	t.Run("Success - Revoke All Sessions Except Current", func(t *testing.T) {
		session1 := models.UserSession{
			UserID:    user.ID,
			TokenHash: "session1_hash",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		testutil.TestDB.Create(&session1)

		session2 := models.UserSession{
			UserID:    user.ID,
			TokenHash: "session2_hash",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		testutil.TestDB.Create(&session2)

		err := userSessionService.RevokeAllSessions(context.Background(), user.ID, &session1.ID)
		require.NoError(t, err)

		activeSessions, err := userSessionService.GetActiveSessions(context.Background(), user.ID)
		require.NoError(t, err)
		assert.Len(t, activeSessions, 1)
		assert.Equal(t, session1.ID, activeSessions[0].ID)
	})

	t.Run("Success - Get Session Stats", func(t *testing.T) {
		stats, err := userSessionService.GetSessionStats(context.Background(), user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, stats.UserID)
		assert.True(t, stats.TotalSessions >= 5)
		assert.Equal(t, int64(1), stats.ActiveSessions)
		assert.True(t, stats.ExpiredSessions >= 1)
		assert.True(t, stats.RevokedSessions >= 3)
	})

	t.Run("Success - Cleanup Expired Sessions", func(t *testing.T) {
		cleaned, err := userSessionService.CleanupExpiredSessions(context.Background())
		require.NoError(t, err)
		assert.True(t, cleaned > 0)
	})
}
