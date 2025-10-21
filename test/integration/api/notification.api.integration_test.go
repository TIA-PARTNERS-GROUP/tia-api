package main
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings" 
	"testing"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants" 
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)
func TestNotificationAPI_Integration_Lifecycle(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()
	
	userA, tokenA := CreateTestUserAndLogin(t, router, "notify-user-a@test.com", "ValidPass123!")
	userB, _ := CreateTestUserAndLogin(t, router, "notify-user-b@test.com", "ValidPass123!") 
	
	constUserBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.UsersBase
	constNotifyBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.NotifyBase
	
	userANotificationsSubPath := strings.Replace(constants.AppRoutes.UserNotifications, ":id", fmt.Sprintf("%d", userA.ID), 1)
	userANotificationsURL := constUserBase + userANotificationsSubPath
	var createdNotification ports.NotificationResponse
	t.Run("Create Notification (System to User A)", func(t *testing.T) {
		createDTO := ports.CreateNotificationInput{
			
			ReceiverUserID:   userA.ID,
			NotificationType: "system",
			Title:            "Welcome!",
			Message:          "Welcome to the platform.",
		}
		body, _ := json.Marshal(createDTO)
		
		req, _ := http.NewRequest(http.MethodPost, constNotifyBase, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenA) 
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		json.Unmarshal(w.Body.Bytes(), &createdNotification)
		assert.NotZero(t, createdNotification.ID)
		assert.Equal(t, userA.ID, createdNotification.Receiver.ID)
		assert.False(t, createdNotification.Read)
		assert.Nil(t, createdNotification.Sender) 
	})
	t.Run("Create Notification (User B to User A)", func(t *testing.T) {
		senderID := userB.ID 
		createDTO := ports.CreateNotificationInput{
			SenderUserID:     &senderID,
			ReceiverUserID:   userA.ID,
			NotificationType: "message",
			Title:            "Hello from B",
			Message:          "Just saying hi.",
		}
		body, _ := json.Marshal(createDTO)
		
		req, _ := http.NewRequest(http.MethodPost, constNotifyBase, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenA)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var notificationFromB ports.NotificationResponse
		json.Unmarshal(w.Body.Bytes(), &notificationFromB)
		assert.NotZero(t, notificationFromB.ID)
		assert.Equal(t, userA.ID, notificationFromB.Receiver.ID)
		assert.NotNil(t, notificationFromB.Sender)
		assert.Equal(t, userB.ID, notificationFromB.Sender.ID)
	})
	t.Run("Get User A Notifications (Unread)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, userANotificationsURL+"?read=false", nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var notifications []ports.NotificationResponse
		json.Unmarshal(w.Body.Bytes(), &notifications)
		assert.Len(t, notifications, 2) 
	})
	t.Run("Fail to Get User B Notifications as User A", func(t *testing.T) {
		
		userBNotificationsSubPath := strings.Replace(constants.AppRoutes.UserNotifications, ":id", fmt.Sprintf("%d", userB.ID), 1)
		userBNotificationsURL := constUserBase + userBNotificationsSubPath
		req, _ := http.NewRequest(http.MethodGet, userBNotificationsURL, nil)
		req.Header.Set("Authorization", "Bearer "+tokenA) 
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code) 
	})
	t.Run("Mark One Notification As Read", func(t *testing.T) {
		
		markReadSubPath := strings.Replace(constants.AppRoutes.UserNotifyReadOne, ":notificationID", fmt.Sprintf("%d", createdNotification.ID), 1)
		markReadURL := userANotificationsURL + markReadSubPath
		req, _ := http.NewRequest(http.MethodPatch, markReadURL, nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var updatedNotification ports.NotificationResponse
		json.Unmarshal(w.Body.Bytes(), &updatedNotification)
		assert.True(t, updatedNotification.Read)
	})
	t.Run("Get User A Notifications (1 Read, 1 Unread)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, userANotificationsURL, nil) 
		req.Header.Set("Authorization", "Bearer "+tokenA)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var notifications []ports.NotificationResponse
		json.Unmarshal(w.Body.Bytes(), &notifications)
		assert.Len(t, notifications, 2)
		
		req, _ = http.NewRequest(http.MethodGet, userANotificationsURL+"?read=true", nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		json.Unmarshal(w.Body.Bytes(), &notifications)
		assert.Len(t, notifications, 1)
		
		req, _ = http.NewRequest(http.MethodGet, userANotificationsURL+"?read=false", nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		json.Unmarshal(w.Body.Bytes(), &notifications)
		assert.Len(t, notifications, 1)
	})
	t.Run("Mark All As Read", func(t *testing.T) {
		
		markAllURL := userANotificationsURL + constants.AppRoutes.UserNotifyReadAll
		req, _ := http.NewRequest(http.MethodPatch, markAllURL, nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var result map[string]int64
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, int64(1), result["marked_as_read_count"]) 
	})
	t.Run("Verify All Are Read", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, userANotificationsURL+"?read=false", nil) 
		req.Header.Set("Authorization", "Bearer "+tokenA)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var notifications []ports.NotificationResponse
		json.Unmarshal(w.Body.Bytes(), &notifications)
		assert.Len(t, notifications, 0) 
	})
	t.Run("Delete Notification", func(t *testing.T) {
		
		deleteSubPath := strings.Replace(constants.AppRoutes.ParamNotificationID, ":notificationID", fmt.Sprintf("%d", createdNotification.ID), 1)
		deleteURL := userANotificationsURL + deleteSubPath
		req, _ := http.NewRequest(http.MethodDelete, deleteURL, nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
	t.Run("Verify Deletion", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, userANotificationsURL, nil) 
		req.Header.Set("Authorization", "Bearer "+tokenA)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var notifications []ports.NotificationResponse
		json.Unmarshal(w.Body.Bytes(), &notifications)
		assert.Len(t, notifications, 1) 
	})
	t.Run("Fail to Delete Notification As Wrong User", func(t *testing.T) {
		
		req, _ := http.NewRequest(http.MethodGet, userANotificationsURL, nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		var notifications []ports.NotificationResponse
		json.Unmarshal(w.Body.Bytes(), &notifications)
		remainingNotificationID := notifications[0].ID
		
		userBNotificationsSubPath := strings.Replace(constants.AppRoutes.UserNotifications, ":id", fmt.Sprintf("%d", userB.ID), 1)
		userBNotificationsURL := constUserBase + userBNotificationsSubPath
		
		deleteSubPath := strings.Replace(constants.AppRoutes.ParamNotificationID, ":notificationID", fmt.Sprintf("%d", remainingNotificationID), 1)
		deleteURLWrongUser := userBNotificationsURL + deleteSubPath
		req, _ = http.NewRequest(http.MethodDelete, deleteURLWrongUser, nil)
		req.Header.Set("Authorization", "Bearer "+tokenA) 
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
