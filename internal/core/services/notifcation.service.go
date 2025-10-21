package services
import (
	"context"
	"errors"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)
type NotificationService struct {
	db *gorm.DB
}
func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}
func (s *NotificationService) CreateNotification(ctx context.Context, data ports.CreateNotificationInput) (*models.Notification, error) {
	if err := s.db.WithContext(ctx).Select("id").First(&models.User{}, data.ReceiverUserID).Error; err != nil {
		return nil, ports.ErrReceiverNotFound
	}
	notification := models.Notification{
		SenderUserID:      data.SenderUserID,
		ReceiverUserID:    data.ReceiverUserID,
		NotificationType:  data.NotificationType,
		Title:             data.Title,
		Message:           data.Message,
		RelatedEntityType: data.RelatedEntityType,
		RelatedEntityID:   data.RelatedEntityID,
		ActionURL:         data.ActionURL,
		Read:              false,
	}
	if err := s.db.WithContext(ctx).Create(&notification).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return s.GetNotificationByID(ctx, notification.ID)
}
func (s *NotificationService) GetNotificationByID(ctx context.Context, id uint) (*models.Notification, error) {
	var notification models.Notification
	err := s.db.WithContext(ctx).
		Preload("ReceiverUser").
		Preload("SenderUser").
		First(&notification, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrNotificationNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &notification, nil
}
func (s *NotificationService) GetNotificationsForUser(ctx context.Context, userID uint, read *bool) ([]models.Notification, error) {
	var notifications []models.Notification
	query := s.db.WithContext(ctx).
		Preload("SenderUser").
		Where("receiver_user_id = ?", userID).
		Order("created_at desc")
	if read != nil {
		query = query.Where("`read` = ?", *read)
	}
	if err := query.Find(&notifications).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return notifications, nil
}
func (s *NotificationService) MarkAsRead(ctx context.Context, id, userID uint) (*models.Notification, error) {
	notification, err := s.GetNotificationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if notification.ReceiverUserID != userID {
		return nil, ports.ErrNotificationNotFound
	}
	if err := s.db.WithContext(ctx).Model(notification).Update("`read`", true).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	notification.Read = true
	return notification, nil
}
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uint) (int64, error) {
	result := s.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("receiver_user_id = ? AND `read` = ?", userID, false).
		Update("`read`", true)
	if result.Error != nil {
		return 0, ports.ErrDatabase
	}
	return result.RowsAffected, nil
}
func (s *NotificationService) DeleteNotification(ctx context.Context, id, userID uint) error {
	result := s.db.WithContext(ctx).
		Delete(&models.Notification{}, "id = ? AND receiver_user_id = ?", id, userID)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrNotificationNotFound
	}
	return nil
}
