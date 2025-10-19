package ports

import (
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

type IdeaVoteType string

const (
	IdeaVoteUp   IdeaVoteType = "up"
	IdeaVoteDown IdeaVoteType = "down"
)

type CreateIdeaVoteInput struct {
	VoterUserID uint         `json:"voter_user_id" validate:"required"`
	IdeaID      uint         `json:"idea_id" validate:"required"`
	VoteType    IdeaVoteType `json:"vote_type" validate:"required,oneof=up down"`
}

type UpdateIdeaVoteInput struct {
	VoteType IdeaVoteType `json:"vote_type" validate:"required,oneof=up down"`
}

type IdeaVoteStats struct {
	IdeaID    uint `json:"idea_id"`
	Upvotes   int  `json:"upvotes"`
	Downvotes int  `json:"downvotes"`
	Total     int  `json:"total_votes"`
	Score     int  `json:"score"`
}

type IdeaVoteResponse struct {
	VoterUserID uint         `json:"voter_user_id"`
	IdeaID      uint         `json:"idea_id"`
	VoteType    IdeaVoteType `json:"vote_type"`
	CreatedAt   time.Time    `json:"created_at"`

	// Relationships
	VoterUser UserResponse `json:"voter_user"`
	Idea      IdeaResponse `json:"idea"`
}

type IdeaVoteStatsResponse struct {
	IdeaID    uint `json:"idea_id"`
	Upvotes   int  `json:"upvotes"`
	Downvotes int  `json:"downvotes"`
	Total     int  `json:"total_votes"`
	Score     int  `json:"score"`
}

type IdeaVotesResponse struct {
	Votes []IdeaVoteResponse `json:"votes"`
	Count int                `json:"count"`
}

func MapToIdeaVoteResponse(iv *models.IdeaVote) IdeaVoteResponse {
	return IdeaVoteResponse{
		VoterUserID: iv.VoterUserID,
		IdeaID:      iv.IdeaID,
		VoteType:    IdeaVoteType(iv.VoteType),
		CreatedAt:   iv.CreatedAt,
		VoterUser:   MapUserToResponse(&iv.VoterUser),
		Idea:        MapIdeaToResponse(&iv.Idea),
	}
}

func MapToIdeaVotesResponse(votes []models.IdeaVote) IdeaVotesResponse {
	voteResponses := make([]IdeaVoteResponse, len(votes))
	for i, vote := range votes {
		voteResponses[i] = MapToIdeaVoteResponse(&vote)
	}

	return IdeaVotesResponse{
		Votes: voteResponses,
		Count: len(voteResponses),
	}
}

func MapToIdeaVoteStats(stats *IdeaVoteStats) IdeaVoteStatsResponse {
	return IdeaVoteStatsResponse{
		IdeaID:    stats.IdeaID,
		Upvotes:   stats.Upvotes,
		Downvotes: stats.Downvotes,
		Total:     stats.Total,
		Score:     stats.Score,
	}
}
