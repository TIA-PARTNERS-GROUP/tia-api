package services

import (
	"context"
	"errors"
	"strings"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)

type DailyActivityService struct {
	db *gorm.DB
}

func NewDailyActivityService(db *gorm.DB) *DailyActivityService {
	return &DailyActivityService{db: db}
}

func (s *DailyActivityService) CreateDailyActivity(ctx context.Context, data ports.CreateDailyActivityInput) (*models.DailyActivity, error) {
	activity := models.DailyActivity{
		Name:        data.Name,
		Description: data.Description,
	}

	if err := s.db.WithContext(ctx).Create(&activity).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ports.ErrActivityNameExists
		}
		return nil, ports.ErrDatabase
	}
	return &activity, nil
}

func (s *DailyActivityService) GetDailyActivityByID(ctx context.Context, id uint) (*models.DailyActivity, error) {
	var activity models.DailyActivity
	if err := s.db.WithContext(ctx).First(&activity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrDailyActivityNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &activity, nil
}

func (s *DailyActivityService) GetAllDailyActivities(ctx context.Context) ([]models.DailyActivity, error) {
	var activities []models.DailyActivity
	if err := s.db.WithContext(ctx).Order("name asc").Find(&activities).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return activities, nil
}

func (s *DailyActivityService) EnrolUserInActivity(ctx context.Context, data ports.EnrolInActivityInput) (*models.DailyActivityEnrolment, error) {
	enrolment := models.DailyActivityEnrolment{
		DailyActivityID: data.DailyActivityID,
		UserID:          data.UserID,
	}

	if err := s.db.WithContext(ctx).Create(&enrolment).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ports.ErrAlreadyEnrolled
		}
		if strings.Contains(err.Error(), "FOREIGN KEY") {
			return nil, ports.ErrProjectOrUserNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &enrolment, nil
}

func (s *DailyActivityService) WithdrawUserFromActivity(ctx context.Context, activityID, userID uint) error {
	result := s.db.WithContext(ctx).Delete(&models.DailyActivityEnrolment{}, "daily_activity_id = ? AND user_id = ?", activityID, userID)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrEnrolmentNotFound
	}
	return nil
}
