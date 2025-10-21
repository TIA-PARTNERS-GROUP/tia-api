package services
import (
	"context"
	"errors"
	"strings"
	"time"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)
type SubscriptionService struct {
	db *gorm.DB
}
func NewSubscriptionService(db *gorm.DB) *SubscriptionService {
	return &SubscriptionService{db: db}
}
func (s *SubscriptionService) CreateSubscription(ctx context.Context, data ports.CreateSubscriptionInput) (*models.Subscription, error) {
	subscription := models.Subscription{
		Name:        data.Name,
		Price:       data.Price,
		ValidDays:   data.ValidDays,
		ValidMonths: data.ValidMonths,
	}
	if err := s.db.WithContext(ctx).Create(&subscription).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ports.ErrSubscriptionNameExists
		}
		return nil, ports.ErrDatabase
	}
	return &subscription, nil
}
func (s *SubscriptionService) GetSubscriptionByID(ctx context.Context, id uint) (*models.Subscription, error) {
	var sub models.Subscription
	if err := s.db.WithContext(ctx).First(&sub, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrSubscriptionNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &sub, nil
}
func (s *SubscriptionService) SubscribeUser(ctx context.Context, data ports.UserSubscribeInput) (*models.UserSubscription, error) {
	plan, err := s.GetSubscriptionByID(ctx, data.SubscriptionID)
	if err != nil {
		return nil, err
	}
	dateFrom := time.Now()
	dateTo := dateFrom
	if plan.ValidDays != nil {
		dateTo = dateFrom.AddDate(0, 0, *plan.ValidDays)
	}
	if plan.ValidMonths != nil {
		dateTo = dateFrom.AddDate(0, *plan.ValidMonths, 0)
	}
	userSub := models.UserSubscription{
		UserID:         data.UserID,
		SubscriptionID: data.SubscriptionID,
		DateFrom:       dateFrom,
		DateTo:         dateTo,
		IsTrial:        false,
	}
	if err := s.db.WithContext(ctx).Create(&userSub).Error; err != nil {
		if strings.Contains(err.Error(), "FOREIGN KEY") {
			return nil, ports.ErrUserNotFound
		}
		return nil, ports.ErrDatabase
	}
	return s.GetUserSubscription(ctx, userSub.ID)
}
func (s *SubscriptionService) GetUserSubscription(ctx context.Context, userSubID uint) (*models.UserSubscription, error) {
	var userSub models.UserSubscription
	err := s.db.WithContext(ctx).
		Preload("User").
		Preload("Subscription").
		First(&userSub, userSubID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrUserSubscriptionNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &userSub, nil
}
