package services
import (
	"context"
	"errors"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)
type EventService struct {
	db *gorm.DB
}
func NewEventService(db *gorm.DB) *EventService {
	return &EventService{db: db}
}
func (s *EventService) CreateEvent(ctx context.Context, data ports.CreateEventInput) (*models.Event, error) {
	event := models.Event{
		EventType: data.EventType,
		Payload:   data.Payload,
		UserID:    data.UserID,
	}
	if err := s.db.WithContext(ctx).Create(&event).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return &event, nil
}
func (s *EventService) GetEventByID(ctx context.Context, id uint) (*models.Event, error) {
	var event models.Event
	err := s.db.WithContext(ctx).Preload("User").First(&event, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrEventNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &event, nil
}
func (s *EventService) GetEvents(ctx context.Context, filters ports.EventsFilter) ([]models.Event, error) {
	var events []models.Event
	query := s.db.WithContext(ctx).Preload("User").Order("timestamp desc")
	if filters.EventType != nil {
		query = query.Where("event_type = ?", *filters.EventType)
	}
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if err := query.Find(&events).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return events, nil
}
