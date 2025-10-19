package ports

import (
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

type CreateIdeaInput struct {
	SubmittedByUserID uint   `json:"submitted_by_user_id" validate:"required"`
	Title             string `json:"title" validate:"required,min=2,max=200"`
	Content           string `json:"content" validate:"required"`
}

type UpdateIdeaInput struct {
	Title   *string `json:"title" validate:"omitempty,min=2,max=200"`
	Content *string `json:"content"`
}

type UpdateIdeaStatusInput struct {
	Status models.IdeaStatus `json:"status" validate:"required"`
}

type VoteInput struct {
	VoterUserID uint                `json:"voter_user_id" validate:"required"`
	VoteType    models.IdeaVoteType `json:"vote_type" validate:"required"`
}

type IdeaResponse struct {
	ID              uint              `json:"id"`
	Title           string            `json:"title"`
	Content         string            `json:"content"`
	Status          models.IdeaStatus `json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	SubmittedByUser UserResponse      `json:"submitted_by_user"`
	Upvotes         int64             `json:"upvotes"`
	Downvotes       int64             `json:"downvotes"`
}

func MapIdeaToResponse(idea *models.Idea) IdeaResponse {
	var upvotes, downvotes int64
	for _, vote := range idea.IdeaVotes {
		if vote.VoteType == models.IdeaVoteUp {
			upvotes++
		} else if vote.VoteType == models.IdeaVoteDown {
			downvotes++
		}
	}

	resp := IdeaResponse{
		ID:        idea.ID,
		Title:     idea.Title,
		Content:   idea.Content,
		Status:    idea.Status,
		CreatedAt: idea.CreatedAt,
		UpdatedAt: idea.UpdatedAt,
		Upvotes:   upvotes,
		Downvotes: downvotes,
	}

	if idea.SubmittedByUser.ID != 0 {
		resp.SubmittedByUser = MapUserToResponse(&idea.SubmittedByUser)
	}

	return resp
}
