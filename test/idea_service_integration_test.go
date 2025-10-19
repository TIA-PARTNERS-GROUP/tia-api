package main

import (
	"context"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/stretchr/testify/assert"
)

func TestIdeaService_Integration_CreateAndGet(t *testing.T) {
	cleanupTestDB(t)
	ideaService := services.NewIdeaService(testDB)

	submitter := models.User{FirstName: "Idea", LoginEmail: "idea@user.com", Active: true}
	testDB.Create(&submitter)

	createDTO := ports.CreateIdeaInput{
		SubmittedByUserID: submitter.ID,
		Title:             "New API Feature",
		Content:           "We should add a new endpoint for ideas.",
	}
	createdIdea, err := ideaService.CreateIdea(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdIdea)
	assert.Equal(t, "New API Feature", createdIdea.Title)
	assert.Equal(t, models.IdeaStatusOpen, createdIdea.Status)

	fetchedIdea, err := ideaService.GetIdeaByID(context.Background(), createdIdea.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedIdea)
	assert.Equal(t, submitter.ID, fetchedIdea.SubmittedByUser.ID)
	assert.Equal(t, "Idea", fetchedIdea.SubmittedByUser.FirstName)
}

func TestIdeaService_Integration_Voting(t *testing.T) {
	cleanupTestDB(t)
	ideaService := services.NewIdeaService(testDB)

	submitter := models.User{FirstName: "Submitter", LoginEmail: "submit@vote.com", Active: true}
	testDB.Create(&submitter)
	voter1 := models.User{FirstName: "Voter1", LoginEmail: "voter1@vote.com", Active: true}
	testDB.Create(&voter1)
	voter2 := models.User{FirstName: "Voter2", LoginEmail: "voter2@vote.com", Active: true}
	testDB.Create(&voter2)

	idea := models.Idea{Title: "Voting Idea", Content: "Content", SubmittedByUserID: submitter.ID}
	testDB.Create(&idea)

	_, err := ideaService.VoteOnIdea(context.Background(), idea.ID, ports.VoteInput{
		VoterUserID: voter1.ID,
		VoteType:    models.IdeaVoteUp,
	})
	assert.NoError(t, err)

	_, err = ideaService.VoteOnIdea(context.Background(), idea.ID, ports.VoteInput{
		VoterUserID: voter2.ID,
		VoteType:    models.IdeaVoteDown,
	})
	assert.NoError(t, err)

	fetchedIdea, _ := ideaService.GetIdeaByID(context.Background(), idea.ID)
	ideaResp := ports.MapIdeaToResponse(fetchedIdea)
	assert.Equal(t, int64(1), ideaResp.Upvotes)
	assert.Equal(t, int64(1), ideaResp.Downvotes)

	_, err = ideaService.VoteOnIdea(context.Background(), idea.ID, ports.VoteInput{
		VoterUserID: voter1.ID,
		VoteType:    models.IdeaVoteDown,
	})
	assert.NoError(t, err)

	fetchedIdea, _ = ideaService.GetIdeaByID(context.Background(), idea.ID)
	ideaResp = ports.MapIdeaToResponse(fetchedIdea)
	assert.Equal(t, int64(0), ideaResp.Upvotes)
	assert.Equal(t, int64(2), ideaResp.Downvotes)
}

func TestIdeaService_Integration_UpdateStatus(t *testing.T) {
	cleanupTestDB(t)
	ideaService := services.NewIdeaService(testDB)

	submitter := models.User{FirstName: "Status", LoginEmail: "status@user.com", Active: true}
	testDB.Create(&submitter)
	idea := models.Idea{Title: "Status Idea", Content: "Content", SubmittedByUserID: submitter.ID, Status: models.IdeaStatusOpen}
	testDB.Create(&idea)

	updatedIdea, err := ideaService.UpdateIdeaStatus(context.Background(), idea.ID, models.IdeaStatusInProgress)
	assert.NoError(t, err)
	assert.Equal(t, models.IdeaStatusInProgress, updatedIdea.Status)
}
