package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
)

func TestIdeaVoteService_Integration_Voting(t *testing.T) {
	cleanupTestDB(t)

	ideaVoteService := services.NewIdeaVoteService(testDB)
	ideaService := services.NewIdeaService(testDB)

	author := models.User{FirstName: "Author", LoginEmail: "author@test.com", Active: true}
	voter := models.User{FirstName: "Voter", LoginEmail: "voter@test.com", Active: true}
	testDB.Create(&author)
	testDB.Create(&voter)

	ideaInput := ports.CreateIdeaInput{
		SubmittedByUserID: author.ID,
		Title:             "Test Idea",
		Content:           "Test Description",
	}
	idea, err := ideaService.CreateIdea(context.Background(), ideaInput)
	require.NoError(t, err)

	t.Run("Success - Create Upvote", func(t *testing.T) {
		input := ports.CreateIdeaVoteInput{
			VoterUserID: voter.ID,
			IdeaID:      idea.ID,
			VoteType:    ports.IdeaVoteUp,
		}

		vote, err := ideaVoteService.CreateIdeaVote(context.Background(), input)
		require.NoError(t, err)
		assert.Equal(t, voter.ID, vote.VoterUserID)
		assert.Equal(t, idea.ID, vote.IdeaID)
		assert.Equal(t, models.IdeaVoteUp, vote.VoteType)
		assert.Equal(t, "Voter", vote.VoterUser.FirstName)
		assert.Equal(t, "Test Idea", vote.Idea.Title)
	})

	t.Run("Failure - Duplicate Vote", func(t *testing.T) {
		input := ports.CreateIdeaVoteInput{
			VoterUserID: voter.ID,
			IdeaID:      idea.ID,
			VoteType:    ports.IdeaVoteDown,
		}

		_, err := ideaVoteService.CreateIdeaVote(context.Background(), input)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrIdeaVoteAlreadyExists, err)
	})

	t.Run("Failure - Vote Own Idea", func(t *testing.T) {
		input := ports.CreateIdeaVoteInput{
			VoterUserID: author.ID,
			IdeaID:      idea.ID,
			VoteType:    ports.IdeaVoteUp,
		}

		_, err := ideaVoteService.CreateIdeaVote(context.Background(), input)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrCannotVoteOwnIdea, err)
	})

	t.Run("Success - Update Vote", func(t *testing.T) {
		updateInput := ports.UpdateIdeaVoteInput{
			VoteType: ports.IdeaVoteDown,
		}

		vote, err := ideaVoteService.UpdateIdeaVote(context.Background(), voter.ID, idea.ID, updateInput)
		require.NoError(t, err)
		assert.Equal(t, models.IdeaVoteDown, vote.VoteType)
	})

	t.Run("Success - Get Vote Stats", func(t *testing.T) {
		stats, err := ideaVoteService.GetIdeaVoteStats(context.Background(), idea.ID)
		require.NoError(t, err)
		assert.Equal(t, idea.ID, stats.IdeaID)
		assert.Equal(t, 0, stats.Upvotes)
		assert.Equal(t, 1, stats.Downvotes)
		assert.Equal(t, 1, stats.Total)
		assert.Equal(t, -1, stats.Score)
	})

	t.Run("Success - Get Votes by Idea", func(t *testing.T) {
		votes, err := ideaVoteService.GetIdeaVotesByIdea(context.Background(), idea.ID)
		require.NoError(t, err)
		assert.Len(t, votes, 1)
		assert.Equal(t, voter.ID, votes[0].VoterUserID)
	})

	t.Run("Success - Get Votes by User", func(t *testing.T) {
		votes, err := ideaVoteService.GetIdeaVotesByUser(context.Background(), voter.ID)
		require.NoError(t, err)
		assert.Len(t, votes, 1)
		assert.Equal(t, idea.ID, votes[0].IdeaID)
	})

	t.Run("Success - Remove Vote", func(t *testing.T) {
		err := ideaVoteService.DeleteIdeaVote(context.Background(), voter.ID, idea.ID)
		require.NoError(t, err)

		_, err = ideaVoteService.GetIdeaVoteByID(context.Background(), voter.ID, idea.ID)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrIdeaVoteNotFound, err)
	})
}

func TestIdeaVoteService_Integration_Validation(t *testing.T) {
	cleanupTestDB(t)

	ideaVoteService := services.NewIdeaVoteService(testDB)

	t.Run("Failure - Non-existent User", func(t *testing.T) {
		input := ports.CreateIdeaVoteInput{
			VoterUserID: 999,
			IdeaID:      1,
			VoteType:    ports.IdeaVoteUp,
		}

		_, err := ideaVoteService.CreateIdeaVote(context.Background(), input)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrUserNotFound, err)
	})

	t.Run("Failure - Non-existent Idea", func(t *testing.T) {
		user := models.User{FirstName: "Test", LoginEmail: "test@test.com", Active: true}
		testDB.Create(&user)

		input := ports.CreateIdeaVoteInput{
			VoterUserID: user.ID,
			IdeaID:      999,
			VoteType:    ports.IdeaVoteUp,
		}

		_, err := ideaVoteService.CreateIdeaVote(context.Background(), input)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrIdeaNotFound, err)
	})
}
