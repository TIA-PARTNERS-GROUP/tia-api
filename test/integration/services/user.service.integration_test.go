package main

import (
	"context"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/TIA-PARTNERS-GROUP/tia-api/pkg/utils"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserService_Integration_CreateUser(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)

	userService := services.NewUserService(testutil.TestDB)

	createDTO := ports.UserCreationSchema{
		FirstName:  "Jane",
		LastName:   "Doe",
		LoginEmail: "jane.doe@example.com",
		Password:   "ValidPassword123!",
	}

	user, err := userService.CreateUser(context.Background(), createDTO)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotZero(t, user.ID)

	var foundUser models.User
	testutil.TestDB.First(&foundUser, user.ID)
	assert.Equal(t, "jane.doe@example.com", foundUser.LoginEmail)
}

func TestUserService_Integration_FindUserByID(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	userService := services.NewUserService(testutil.TestDB)

	seededUser := models.User{FirstName: "Segun", LoginEmail: "segun@example.com", Active: true}
	testutil.TestDB.Create(&seededUser)

	t.Run("Success - User Found", func(t *testing.T) {
		user, err := userService.FindUserByID(context.Background(), seededUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "Segun", user.FirstName)
	})

	t.Run("Failure - User Not Found", func(t *testing.T) {
		_, err := userService.FindUserByID(context.Background(), 9999)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrUserNotFound, err)
	})
}

func TestUserService_Integration_UpdateUser(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	userService := services.NewUserService(testutil.TestDB)

	seededUser := models.User{FirstName: "Original", LoginEmail: "update@me.com", Active: true}
	testutil.TestDB.Create(&seededUser)

	newName := "Updated"
	updateDTO := ports.UserUpdateSchema{FirstName: &newName}

	updatedUser, err := userService.UpdateUser(context.Background(), seededUser.ID, updateDTO)

	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)

	var foundUser models.User
	testutil.TestDB.First(&foundUser, seededUser.ID)
	assert.Equal(t, "Updated", foundUser.FirstName)
}

func TestUserService_Integration_ChangePassword(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	userService := services.NewUserService(testutil.TestDB)

	oldPassword := "oldPassword123!"
	newPassword := "newValidPassword456!"
	hashedOldPassword, _ := utils.HashPassword(oldPassword)
	seededUser := models.User{FirstName: "Pass", LoginEmail: "pass@change.com", PasswordHash: &hashedOldPassword, Active: true}
	testutil.TestDB.Create(&seededUser)

	t.Run("Success", func(t *testing.T) {
		err := userService.ChangePassword(context.Background(), seededUser.ID, oldPassword, newPassword)
		assert.NoError(t, err)

		var foundUser models.User
		testutil.TestDB.First(&foundUser, seededUser.ID)
		assert.NoError(t, utils.VerifyPassword(newPassword, *foundUser.PasswordHash))
	})

	t.Run("Failure - Incorrect current password", func(t *testing.T) {
		err := userService.ChangePassword(context.Background(), seededUser.ID, "wrong-password", newPassword)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrIncorrectPassword, err)
	})
}

func TestUserService_Integration_DeleteUser(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	userService := services.NewUserService(testutil.TestDB)

	seededUser := models.User{FirstName: "ToDelete", LoginEmail: "delete@me.com", Active: true}
	testutil.TestDB.Create(&seededUser)

	err := userService.DeleteUser(context.Background(), seededUser.ID)
	assert.NoError(t, err)

	var foundUser models.User
	result := testutil.TestDB.First(&foundUser, seededUser.ID)
	assert.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
}

func TestUserService_Integration_FindAllUsers(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	userService := services.NewUserService(testutil.TestDB)

	testutil.TestDB.Create(&models.User{FirstName: "User1", LoginEmail: "user1@test.com", Active: true})
	testutil.TestDB.Create(&models.User{FirstName: "User2", LoginEmail: "user2@test.com", Active: true})

	users, err := userService.FindAllUsers(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUserService_Integration_DeactivateUser(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	userService := services.NewUserService(testutil.TestDB)

	seededUser := models.User{FirstName: "ToDeactivate", LoginEmail: "deactivate@me.com", Active: true}
	testutil.TestDB.Create(&seededUser)

	user, err := userService.DeactivateUser(context.Background(), seededUser.ID)

	assert.NoError(t, err)
	assert.False(t, user.Active)

	var foundUser models.User
	testutil.TestDB.First(&foundUser, seededUser.ID)
	assert.False(t, foundUser.Active)
}
