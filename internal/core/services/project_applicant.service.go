package services

import (
	"context"
	"strings"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)

type ProjectApplicantService struct {
	db *gorm.DB
}

func NewProjectApplicantService(db *gorm.DB) *ProjectApplicantService {
	return &ProjectApplicantService{db: db}
}
func (s *ProjectApplicantService) ApplyToProject(ctx context.Context, data ports.ApplyToProjectInput) (*models.ProjectApplicant, error) {
	application := models.ProjectApplicant{
		ProjectID: data.ProjectID,
		UserID:    data.UserID,
	}
	if err := s.db.WithContext(ctx).Create(&application).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ports.ErrAlreadyApplied
		}
		if strings.Contains(err.Error(), "FOREIGN KEY") {
			return nil, ports.ErrProjectOrUserNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &application, nil
}
func (s *ProjectApplicantService) WithdrawApplication(ctx context.Context, projectID, userID uint) error {
	result := s.db.WithContext(ctx).Delete(&models.ProjectApplicant{}, "project_id = ? AND user_id = ?", projectID, userID)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrApplicationNotFound
	}
	return nil
}
func (s *ProjectApplicantService) GetApplicantsForProject(ctx context.Context, projectID uint) ([]models.ProjectApplicant, error) {
	var applicants []models.ProjectApplicant
	err := s.db.WithContext(ctx).
		Preload("User").
		Where("project_id = ?", projectID).
		Find(&applicants).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return applicants, nil
}
func (s *ProjectApplicantService) GetApplicationsForUser(ctx context.Context, userID uint) ([]models.ProjectApplicant, error) {
	var applications []models.ProjectApplicant
	err := s.db.WithContext(ctx).
		Preload("Project").
		Where("user_id = ?", userID).
		Find(&applications).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return applications, nil
}
