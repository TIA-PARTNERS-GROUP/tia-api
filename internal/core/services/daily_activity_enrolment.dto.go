package services
import (
	"context"
	"strings"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)
type DailyActivityEnrolmentService struct {
	db *gorm.DB
}
func NewDailyActivityEnrolmentService(db *gorm.DB) *DailyActivityEnrolmentService {
	return &DailyActivityEnrolmentService{db: db}
}
func (s *DailyActivityEnrolmentService) EnrolUser(ctx context.Context, data ports.EnrolmentInput) (*models.DailyActivityEnrolment, error) {
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
func (s *DailyActivityEnrolmentService) WithdrawUser(ctx context.Context, activityID, userID uint) error {
	result := s.db.WithContext(ctx).Delete(&models.DailyActivityEnrolment{}, "daily_activity_id = ? AND user_id = ?", activityID, userID)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrEnrolmentNotFound
	}
	return nil
}
func (s *DailyActivityEnrolmentService) GetEnrolmentsForActivity(ctx context.Context, activityID uint) ([]models.DailyActivityEnrolment, error) {
	var enrolments []models.DailyActivityEnrolment
	err := s.db.WithContext(ctx).
		Preload("User").
		Where("daily_activity_id = ?", activityID).
		Find(&enrolments).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return enrolments, nil
}
func (s *DailyActivityEnrolmentService) GetEnrolmentsForUser(ctx context.Context, userID uint) ([]models.DailyActivityEnrolment, error) {
	var enrolments []models.DailyActivityEnrolment
	err := s.db.WithContext(ctx).
		Preload("DailyActivity").
		Where("user_id = ?", userID).
		Find(&enrolments).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return enrolments, nil
}
