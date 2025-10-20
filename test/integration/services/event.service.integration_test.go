package main

import (
	"context"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestEventService_Integration_CreateAndGet(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	eventService := services.NewEventService(testutil.TestDB)

	user := models.User{FirstName: "EventUser", LoginEmail: "event@user.com", Active: true}
	testutil.TestDB.Create(&user)

	t.Run("Create event with user", func(t *testing.T) {
		payload := `{"page": "/dashboard", "action": "view"}`
		createDTO := ports.CreateEventInput{
			EventType: "page_view",
			Payload:   datatypes.JSON(payload),
			UserID:    &user.ID,
		}

		createdEvent, err := eventService.CreateEvent(context.Background(), createDTO)
		assert.NoError(t, err)
		assert.NotNil(t, createdEvent)
		assert.Equal(t, "page_view", createdEvent.EventType)
		assert.Equal(t, user.ID, *createdEvent.UserID)

		fetchedEvent, err := eventService.GetEventByID(context.Background(), createdEvent.ID)
		assert.NoError(t, err)
		assert.NotNil(t, fetchedEvent.User)
		assert.Equal(t, user.ID, fetchedEvent.User.ID)
	})

	t.Run("Create system event without user", func(t *testing.T) {
		payload := `{"service": "cleanup_job", "status": "started"}`
		createDTO := ports.CreateEventInput{
			EventType: "system_job",
			Payload:   datatypes.JSON(payload),
			UserID:    nil,
		}
		createdEvent, err := eventService.CreateEvent(context.Background(), createDTO)
		assert.NoError(t, err)
		assert.Nil(t, createdEvent.UserID)

		fetchedEvent, err := eventService.GetEventByID(context.Background(), createdEvent.ID)
		assert.NoError(t, err)
		assert.Nil(t, fetchedEvent.User)
	})
}

func TestEventService_Integration_GetEventsFiltered(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	eventService := services.NewEventService(testutil.TestDB)

	user1 := models.User{FirstName: "EventUser1", LoginEmail: "event1@user.com", Active: true}
	testutil.TestDB.Create(&user1)
	user2 := models.User{FirstName: "EventUser2", LoginEmail: "event2@user.com", Active: true}
	testutil.TestDB.Create(&user2)

	testutil.TestDB.Create(&models.Event{EventType: "login", UserID: &user1.ID, Payload: datatypes.JSON(`{}`)})
	testutil.TestDB.Create(&models.Event{EventType: "login", UserID: &user2.ID, Payload: datatypes.JSON(`{}`)})
	testutil.TestDB.Create(&models.Event{EventType: "logout", UserID: &user1.ID, Payload: datatypes.JSON(`{}`)})
	testutil.TestDB.Create(&models.Event{EventType: "system_error", UserID: nil, Payload: datatypes.JSON(`{}`)})

	t.Run("Filter by EventType", func(t *testing.T) {
		eventType := "login"
		filters := ports.EventsFilter{EventType: &eventType}
		events, err := eventService.GetEvents(context.Background(), filters)
		assert.NoError(t, err)
		assert.Len(t, events, 2)
	})

	t.Run("Filter by UserID", func(t *testing.T) {
		filters := ports.EventsFilter{UserID: &user1.ID}
		events, err := eventService.GetEvents(context.Background(), filters)
		assert.NoError(t, err)
		assert.Len(t, events, 2)
	})

	t.Run("Filter by EventType and UserID", func(t *testing.T) {
		eventType := "logout"
		filters := ports.EventsFilter{EventType: &eventType, UserID: &user1.ID}
		events, err := eventService.GetEvents(context.Background(), filters)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, "logout", events[0].EventType)
	})
}
