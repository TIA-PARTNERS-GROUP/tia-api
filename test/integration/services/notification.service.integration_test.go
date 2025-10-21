package main
import (
	"context"
	"testing"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)
func TestNotificationService_Integration_CreateAndGet(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	notifService := services.NewNotificationService(testutil.TestDB)
	sender := models.User{FirstName: "Sender", LoginEmail: "sender@notif.com", Active: true}
	testutil.TestDB.Create(&sender)
	receiver := models.User{FirstName: "Receiver", LoginEmail: "receiver@notif.com", Active: true}
	testutil.TestDB.Create(&receiver)
	createDTO := ports.CreateNotificationInput{
		SenderUserID:     &sender.ID,
		ReceiverUserID:   receiver.ID,
		NotificationType: "project_invite",
		Title:            "You're invited!",
		Message:          "Join our new project.",
	}
	createdNotif, err := notifService.CreateNotification(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdNotif)
	assert.Equal(t, "You're invited!", createdNotif.Title)
	assert.False(t, createdNotif.Read)
	fetchedNotif, err := notifService.GetNotificationByID(context.Background(), createdNotif.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedNotif)
	assert.NotNil(t, fetchedNotif.SenderUser)
	assert.NotNil(t, fetchedNotif.ReceiverUser)
	assert.Equal(t, sender.ID, fetchedNotif.SenderUser.ID)
	assert.Equal(t, receiver.ID, fetchedNotif.ReceiverUser.ID)
}
func TestNotificationService_Integration_FilteringAndStatus(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	notifService := services.NewNotificationService(testutil.TestDB)
	user := models.User{FirstName: "NotifUser", LoginEmail: "notif@user.com", Active: true}
	result := testutil.TestDB.Create(&user)
	assert.NoError(t, result.Error)
	t.Logf("Created user with ID: %d", user.ID)
	
	notifications := []models.Notification{
		{
			ReceiverUserID:   user.ID,
			NotificationType: "system",
			Title:            "Unread 1",
			Message:          "...",
			Read:             false,
		},
		{
			ReceiverUserID:   user.ID,
			NotificationType: "system",
			Title:            "Unread 2",
			Message:          "...",
			Read:             false,
		},
		{
			ReceiverUserID:   user.ID,
			NotificationType: "system",
			Title:            "Read 1",
			Message:          "...",
			Read:             true,
		},
	}
	for i, notif := range notifications {
		result := testutil.TestDB.Create(&notif)
		assert.NoError(t, result.Error)
		t.Logf("Created notification %d with ID: %d, Read: %t", i+1, notif.ID, notif.Read)
	}
	
	var count int64
	testutil.TestDB.Model(&models.Notification{}).Where("receiver_user_id = ?", user.ID).Count(&count)
	t.Logf("Total notifications in DB for user: %d", count)
	t.Run("Get All Notifications", func(t *testing.T) {
		notifs, err := notifService.GetNotificationsForUser(context.Background(), user.ID, nil)
		t.Logf("Get All - Error: %v, Count: %d", err, len(notifs))
		assert.NoError(t, err)
		assert.Len(t, notifs, 3)
	})
	t.Run("Get Unread Notifications", func(t *testing.T) {
		readStatus := false
		notifs, err := notifService.GetNotificationsForUser(context.Background(), user.ID, &readStatus)
		t.Logf("Get Unread - Error: %v, Count: %d", err, len(notifs))
		assert.NoError(t, err)
		assert.Len(t, notifs, 2)
	})
	t.Run("Mark Single As Read", func(t *testing.T) {
		var unread models.Notification
		result := testutil.TestDB.First(&unread, "receiver_user_id = ? AND `read` = ?", user.ID, false)
		t.Logf("Find unread - Error: %v, Found ID: %d", result.Error, unread.ID)
		assert.NoError(t, result.Error)
		updatedNotif, err := notifService.MarkAsRead(context.Background(), unread.ID, user.ID)
		t.Logf("MarkAsRead - Error: %v, Updated Read: %t", err, updatedNotif.Read)
		assert.NoError(t, err)
		assert.True(t, updatedNotif.Read)
		readStatus := false
		notifs, _ := notifService.GetNotificationsForUser(context.Background(), user.ID, &readStatus)
		t.Logf("After MarkAsRead - Unread count: %d", len(notifs))
		assert.Len(t, notifs, 1)
	})
	t.Run("Mark All As Read", func(t *testing.T) {
		rowsAffected, err := notifService.MarkAllAsRead(context.Background(), user.ID)
		t.Logf("MarkAllAsRead - Error: %v, Rows affected: %d", err, rowsAffected)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), rowsAffected)
		readStatus := false
		notifs, _ := notifService.GetNotificationsForUser(context.Background(), user.ID, &readStatus)
		t.Logf("After MarkAllAsRead - Unread count: %d", len(notifs))
		assert.Len(t, notifs, 0)
	})
}
