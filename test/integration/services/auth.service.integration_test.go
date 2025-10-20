package main

import (
	"context"
	"testing"
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/TIA-PARTNERS-GROUP/tia-api/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Integration_Login(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	authService := services.NewAuthService(testutil.TestDB)

	password := "ValidPassword123!"
	hashedPassword, _ := utils.HashPassword(password)
	seededUser := models.User{
		FirstName:    "Auth",
		LoginEmail:   "auth@test.com",
		PasswordHash: &hashedPassword,
		Active:       true,
	}
	testutil.TestDB.Create(&seededUser)

	t.Run("Success - Valid Credentials", func(t *testing.T) {
		loginDTO := ports.LoginInput{
			LoginEmail: "auth@test.com",
			Password:   password,
		}
		res, err := authService.Login(context.Background(), loginDTO, nil, nil)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, seededUser.ID, res.User.ID)
	})

	t.Run("Failure - Invalid Password", func(t *testing.T) {
		loginDTO := ports.LoginInput{
			LoginEmail: "auth@test.com",
			Password:   "wrong-password",
		}
		_, err := authService.Login(context.Background(), loginDTO, nil, nil)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrInvalidCredentials, err)
	})

	t.Run("Failure - User Not Found", func(t *testing.T) {
		loginDTO := ports.LoginInput{
			LoginEmail: "notfound@test.com",
			Password:   password,
		}
		_, err := authService.Login(context.Background(), loginDTO, nil, nil)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrInvalidCredentials, err)
	})

	t.Run("Failure - Deactivated Account", func(t *testing.T) {
		password := "ValidPassword123!"
		hashedPassword, _ := utils.HashPassword(password)

		deactivatedUser := models.User{
			FirstName:    "Inactive",
			LoginEmail:   "inactive@test.com",
			PasswordHash: &hashedPassword,
		}

		testutil.TestDB.Create(&deactivatedUser)
		testutil.TestDB.Model(&deactivatedUser).Update("active", false)

		loginDTO := ports.LoginInput{
			LoginEmail: "inactive@test.com",
			Password:   password,
		}
		_, err := authService.Login(context.Background(), loginDTO, nil, nil)

		assert.Error(t, err)
		assert.Equal(t, ports.ErrAccountDeactivated, err)
	})
}

func TestAuthService_Integration_Logout(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	authService := services.NewAuthService(testutil.TestDB)
	seededUser := models.User{FirstName: "Logout", LoginEmail: "logout@test.com", Active: true}
	testutil.TestDB.Create(&seededUser)
	activeSession := models.UserSession{UserID: seededUser.ID, TokenHash: "somehash", ExpiresAt: time.Now().Add(time.Hour)}
	testutil.TestDB.Create(&activeSession)
	success, err := authService.Logout(context.Background(), activeSession.ID, seededUser.ID)
	assert.NoError(t, err)
	assert.True(t, success)
	var foundSession models.UserSession
	testutil.TestDB.First(&foundSession, activeSession.ID)
	assert.NotNil(t, foundSession.RevokedAt)
}

func TestAuthService_Integration_ValidateToken(t *testing.T) {
	authService := services.NewAuthService(testutil.TestDB)
	seededUser := models.User{FirstName: "Token", LoginEmail: "token@test.com", Active: true}
	testutil.TestDB.Create(&seededUser)

	t.Run("Success - Valid Token and Session", func(t *testing.T) {
		session := models.UserSession{UserID: seededUser.ID}
		testutil.TestDB.Create(&session)
		token, expiry, _ := utils.GenerateToken(seededUser.ID, session.ID, seededUser.LoginEmail)
		session.TokenHash = utils.HashToken(token)
		session.ExpiresAt = expiry
		testutil.TestDB.Save(&session)
		user, sess, err := authService.ValidateToken(context.Background(), token)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotNil(t, sess)
	})

	t.Run("Failure - Invalid Token String", func(t *testing.T) {
		_, _, err := authService.ValidateToken(context.Background(), "this.is.a.bad.token")
		assert.Error(t, err)
		assert.Equal(t, ports.ErrInvalidToken, err)
	})

	t.Run("Failure - Session Not Found", func(t *testing.T) {
		token, _, _ := utils.GenerateToken(seededUser.ID, 9999, seededUser.LoginEmail)
		_, _, err := authService.ValidateToken(context.Background(), token)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrInvalidSession, err)
	})
}
