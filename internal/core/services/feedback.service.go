package services

import (
	"context"
	"errors"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)

type FeedbackService struct {
	db *gorm.DB
}

func NewFeedbackService(db *gorm.DB) *FeedbackService {
	return &FeedbackService{db: db}
}

func (s *FeedbackService) CreateFeedback(ctx context.Context, data ports.CreateFeedbackInput) (*models.Feedback, error) {
	feedback := models.Feedback{
		Name:    data.Name,
		Email:   data.Email,
		Content: data.Content,
	}

	if err := s.db.WithContext(ctx).Create(&feedback).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return &feedback, nil
}

func (s *FeedbackService) GetFeedbackByID(ctx context.Context, id uint) (*models.Feedback, error) {
	var feedback models.Feedback
	if err := s.db.WithContext(ctx).First(&feedback, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrFeedbackNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &feedback, nil
}

func (s *FeedbackService) GetAllFeedback(ctx context.Context) ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	if err := s.db.WithContext(ctx).Order("date_submitted desc").Find(&feedbacks).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return feedbacks, nil
}

func (s *FeedbackService) DeleteFeedback(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Feedback{}, id)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrFeedbackNotFound
	}
	return nil
}
