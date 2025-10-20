package services

import (
	"context"
	"errors"
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)

type UserSubscriptionService struct {
	db *gorm.DB
}

func NewUserSubscriptionService(db *gorm.DB) *UserSubscriptionService {
	return &UserSubscriptionService{db: db}
}

func (s *UserSubscriptionService) CreateUserSubscription(ctx context.Context, data ports.CreateUserSubscriptionInput) (*models.UserSubscription, error) {
	var plan models.Subscription
	if err := s.db.WithContext(ctx).First(&plan, data.SubscriptionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrSubscriptionNotFound
		}
		return nil, ports.ErrDatabase
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
		IsTrial:        data.IsTrial,
	}

	if err := s.db.WithContext(ctx).Create(&userSub).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	return s.GetUserSubscriptionByID(ctx, userSub.ID)
}

func (s *UserSubscriptionService) GetUserSubscriptionByID(ctx context.Context, id uint) (*models.UserSubscription, error) {
	var userSub models.UserSubscription
	err := s.db.WithContext(ctx).
		Preload("User").
		Preload("Subscription").
		First(&userSub, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrUserSubscriptionNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &userSub, nil
}

func (s *UserSubscriptionService) GetSubscriptionsForUser(ctx context.Context, userID uint) ([]models.UserSubscription, error) {
	var userSubs []models.UserSubscription
	err := s.db.WithContext(ctx).
		Preload("Subscription").
		Where("user_id = ? AND date_to > ?", userID, time.Now()).
		Find(&userSubs).Error

	if err != nil {
		return nil, ports.ErrDatabase
	}
	return userSubs, nil
}

func (s *UserSubscriptionService) CancelSubscription(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.UserSubscription{}, id)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrUserSubscriptionNotFound
	}
	return nil
}
