package services

import (
	"context"
	"errors"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IdeaService struct {
	db *gorm.DB
}

func NewIdeaService(db *gorm.DB) *IdeaService {
	return &IdeaService{db: db}
}

func (s *IdeaService) CreateIdea(ctx context.Context, data ports.CreateIdeaInput) (*models.Idea, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, data.SubmittedByUserID).Error; err != nil {
		return nil, ports.ErrIdeaSubmitterNotFound
	}

	idea := models.Idea{
		SubmittedByUserID: data.SubmittedByUserID,
		Title:             data.Title,
		Content:           data.Content,
		Status:            models.IdeaStatusOpen,
	}

	if err := s.db.WithContext(ctx).Create(&idea).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return s.GetIdeaByID(ctx, idea.ID)
}

func (s *IdeaService) GetIdeaByID(ctx context.Context, id uint) (*models.Idea, error) {
	var idea models.Idea
	err := s.db.WithContext(ctx).
		Preload("SubmittedByUser").
		Preload("IdeaVotes").
		First(&idea, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrIdeaNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &idea, nil
}

func (s *IdeaService) GetAllIdeas(ctx context.Context, status *models.IdeaStatus) ([]models.Idea, error) {
	var ideas []models.Idea
	query := s.db.WithContext(ctx).Preload("SubmittedByUser").Preload("IdeaVotes").Order("created_at desc")

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Find(&ideas).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return ideas, nil
}

func (s *IdeaService) UpdateIdeaStatus(ctx context.Context, id uint, status models.IdeaStatus) (*models.Idea, error) {
	if err := s.db.WithContext(ctx).Model(&models.Idea{ID: id}).Update("status", status).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return s.GetIdeaByID(ctx, id)
}

func (s *IdeaService) VoteOnIdea(ctx context.Context, id uint, data ports.VoteInput) (*models.Idea, error) {
	vote := models.IdeaVote{
		IdeaID:      id,
		VoterUserID: data.VoterUserID,
		VoteType:    data.VoteType,
	}

	err := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "idea_id"}, {Name: "voter_user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"vote_type"}),
	}).Create(&vote).Error

	if err != nil {
		return nil, ports.ErrDatabase
	}

	return s.GetIdeaByID(ctx, id)
}

func (s *IdeaService) DeleteIdea(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Idea{}, id)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrIdeaNotFound
	}
	return nil
}
