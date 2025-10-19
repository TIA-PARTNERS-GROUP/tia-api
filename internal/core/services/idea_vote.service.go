package services

import (
	"context"
	"errors"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)

type IdeaVoteService struct {
	db *gorm.DB
}

func NewIdeaVoteService(db *gorm.DB) *IdeaVoteService {
	return &IdeaVoteService{db: db}
}

func (s *IdeaVoteService) CreateIdeaVote(ctx context.Context, data ports.CreateIdeaVoteInput) (*models.IdeaVote, error) {
	// Check if user exists
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, data.VoterUserID).Error; err != nil {
		return nil, ports.ErrUserNotFound
	}

	// Check if idea exists
	var idea models.Idea
	if err := s.db.WithContext(ctx).First(&idea, data.IdeaID).Error; err != nil {
		return nil, ports.ErrIdeaNotFound
	}

	// Check if user is trying to vote on their own idea
	if idea.SubmittedByUserID == data.VoterUserID {
		return nil, ports.ErrCannotVoteOwnIdea
	}

	// Check if vote already exists
	var existingVote models.IdeaVote
	err := s.db.WithContext(ctx).
		Where("idea_id = ? AND voter_user_id = ?", data.IdeaID, data.VoterUserID).
		First(&existingVote).Error

	if err == nil {
		return nil, ports.ErrIdeaVoteAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ports.ErrDatabase
	}

	// Create the vote
	vote := models.IdeaVote{
		IdeaID:      data.IdeaID,
		VoterUserID: data.VoterUserID,
		VoteType:    models.IdeaVoteType(data.VoteType),
	}

	if err := s.db.WithContext(ctx).Create(&vote).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	return s.GetIdeaVoteByID(ctx, data.VoterUserID, data.IdeaID)
}

func (s *IdeaVoteService) GetIdeaVoteByID(ctx context.Context, voterUserID, ideaID uint) (*models.IdeaVote, error) {
	var vote models.IdeaVote
	err := s.db.WithContext(ctx).
		Preload("VoterUser").
		Preload("Idea").
		Where("voter_user_id = ? AND idea_id = ?", voterUserID, ideaID).
		First(&vote).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrIdeaVoteNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &vote, nil
}

func (s *IdeaVoteService) GetIdeaVotesByIdea(ctx context.Context, ideaID uint) ([]models.IdeaVote, error) {
	var votes []models.IdeaVote
	err := s.db.WithContext(ctx).
		Preload("VoterUser").
		Where("idea_id = ?", ideaID).
		Find(&votes).Error

	if err != nil {
		return nil, ports.ErrDatabase
	}
	return votes, nil
}

func (s *IdeaVoteService) GetIdeaVotesByUser(ctx context.Context, userID uint) ([]models.IdeaVote, error) {
	var votes []models.IdeaVote
	err := s.db.WithContext(ctx).
		Preload("Idea").
		Where("voter_user_id = ?", userID).
		Find(&votes).Error

	if err != nil {
		return nil, ports.ErrDatabase
	}
	return votes, nil
}

func (s *IdeaVoteService) UpdateIdeaVote(ctx context.Context, voterUserID, ideaID uint, data ports.UpdateIdeaVoteInput) (*models.IdeaVote, error) {
	var vote models.IdeaVote
	err := s.db.WithContext(ctx).
		Where("voter_user_id = ? AND idea_id = ?", voterUserID, ideaID).
		First(&vote).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrIdeaVoteNotFound
		}
		return nil, ports.ErrDatabase
	}

	vote.VoteType = models.IdeaVoteType(data.VoteType)

	if err := s.db.WithContext(ctx).Save(&vote).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	return s.GetIdeaVoteByID(ctx, voterUserID, ideaID)
}

func (s *IdeaVoteService) DeleteIdeaVote(ctx context.Context, voterUserID, ideaID uint) error {
	result := s.db.WithContext(ctx).
		Where("voter_user_id = ? AND idea_id = ?", voterUserID, ideaID).
		Delete(&models.IdeaVote{})

	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrIdeaVoteNotFound
	}
	return nil
}

func (s *IdeaVoteService) GetIdeaVoteStats(ctx context.Context, ideaID uint) (*ports.IdeaVoteStatsResponse, error) {
	var upvotes, downvotes int64

	// Count upvotes
	if err := s.db.WithContext(ctx).Model(&models.IdeaVote{}).
		Where("idea_id = ? AND vote_type = ?", ideaID, models.IdeaVoteUp).
		Count(&upvotes).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	// Count downvotes
	if err := s.db.WithContext(ctx).Model(&models.IdeaVote{}).
		Where("idea_id = ? AND vote_type = ?", ideaID, models.IdeaVoteDown).
		Count(&downvotes).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	stats := &ports.IdeaVoteStatsResponse{
		IdeaID:    ideaID,
		Upvotes:   int(upvotes),
		Downvotes: int(downvotes),
		Total:     int(upvotes + downvotes),
		Score:     int(upvotes - downvotes),
	}

	return stats, nil
}
