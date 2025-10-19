package services

import (
	"context"
	"errors"
	"strings"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)

type ProjectService struct {
	db *gorm.DB
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{db: db}
}

func (s *ProjectService) CreateProject(ctx context.Context, data ports.CreateProjectInput) (*models.Project, error) {
	var manager models.User
	if err := s.db.WithContext(ctx).First(&manager, data.ManagedByUserID).Error; err != nil {
		return nil, ports.ErrManagerNotFound
	}

	project := models.Project{
		ManagedByUserID: data.ManagedByUserID,
		BusinessID:      data.BusinessID,
		Name:            data.Name,
		Description:     data.Description,
		ProjectStatus:   data.ProjectStatus,
		StartDate:       data.StartDate,
		TargetEndDate:   data.TargetEndDate,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&project).Error; err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				return ports.ErrProjectNameExists
			}
			return err
		}

		member := models.ProjectMember{
			ProjectID: project.ID,
			UserID:    project.ManagedByUserID,
			Role:      models.ProjectMemberRoleManager,
		}
		if err := tx.Create(&member).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetProjectByID(ctx, project.ID)
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id uint) (*models.Project, error) {
	var project models.Project
	err := s.db.WithContext(ctx).
		Preload("ManagingUser").
		Preload("Business").
		Preload("ProjectMembers.User").
		First(&project, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrProjectNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &project, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, id uint, data ports.UpdateProjectInput) (*models.Project, error) {
	_, err := s.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updateData := make(map[string]interface{})

	if data.ManagedByUserID != nil {
		updateData["managed_by_user_id"] = *data.ManagedByUserID
	}
	if data.BusinessID != nil {
		updateData["business_id"] = *data.BusinessID
	}
	if data.Name != nil {
		updateData["name"] = *data.Name
	}
	if data.Description != nil {
		updateData["description"] = *data.Description
	}
	if data.ProjectStatus != nil {
		updateData["project_status"] = *data.ProjectStatus
	}
	if data.StartDate != nil {
		updateData["start_date"] = *data.StartDate
	}
	if data.TargetEndDate != nil {
		updateData["target_end_date"] = *data.TargetEndDate
	}
	if data.ActualEndDate != nil {
		updateData["actual_end_date"] = *data.ActualEndDate
	}

	if len(updateData) == 0 {
		return nil, ports.ErrNoUpdateData
	}

	if len(updateData) == 0 {
		return nil, ports.ErrNoUpdateData
	}

	if err := s.db.WithContext(ctx).Model(&models.Project{ID: id}).Updates(updateData).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	return s.GetProjectByID(ctx, id)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Project{}, id)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrProjectNotFound
	}
	return nil
}

func (s *ProjectService) AddMember(ctx context.Context, projectID uint, data ports.AddMemberInput) (*models.ProjectMember, error) {
	member := models.ProjectMember{
		ProjectID: projectID,
		UserID:    data.UserID,
		Role:      data.Role,
	}

	if err := s.db.WithContext(ctx).Create(&member).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ports.ErrMemberAlreadyExists
		}
		return nil, ports.ErrDatabase
	}

	var createdMember models.ProjectMember
	s.db.WithContext(ctx).Preload("User").First(&createdMember, "project_id = ? AND user_id = ?", projectID, data.UserID)
	return &createdMember, nil
}

func (s *ProjectService) RemoveMember(ctx context.Context, projectID, userID uint) error {
	result := s.db.WithContext(ctx).Delete(&models.ProjectMember{}, "project_id = ? AND user_id = ?", projectID, userID)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrMemberNotFound
	}
	return nil
}
