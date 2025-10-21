package main
import (
	"context"
	"testing"
	"time"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)
func TestPublicationService_Integration_CreateAndGet(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	pubService := services.NewPublicationService(testutil.TestDB)
	author := models.User{FirstName: "Author", LoginEmail: "author@pub.com", Active: true}
	testutil.TestDB.Create(&author)
	published := true
	createDTO := ports.CreatePublicationInput{
		UserID:          author.ID,
		PublicationType: models.PublicationPost,
		Title:           "My First Go Post!",
		Content:         "This is the content.",
		Published:       &published,
	}
	createdPub, err := pubService.CreatePublication(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdPub)
	assert.Equal(t, "My First Go Post!", createdPub.Title)
	assert.Equal(t, "my-first-go-post", createdPub.Slug)
	assert.True(t, createdPub.Published)
	assert.NotNil(t, createdPub.PublishedAt)
	fetchedPub, err := pubService.GetPublicationBySlug(context.Background(), "my-first-go-post")
	assert.NoError(t, err)
	assert.NotNil(t, fetchedPub)
	assert.Equal(t, "My First Go Post!", fetchedPub.Title)
	assert.NotZero(t, fetchedPub.User.ID)
	assert.Equal(t, "Author", fetchedPub.User.FirstName)
}
func TestPublicationService_Integration_UpdatePublication(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	pubService := services.NewPublicationService(testutil.TestDB)
	author := models.User{FirstName: "Author", LoginEmail: "author@pub.com", Active: true}
	testutil.TestDB.Create(&author)
	pub := models.Publication{
		UserID:          author.ID,
		PublicationType: models.PublicationArticle,
		Title:           "Original Title",
		Slug:            "original-title",
		Content:         "Original content.",
		Published:       false,
	}
	testutil.TestDB.Select("*").Create(&pub)
	newTitle := "Updated Awesome Title"
	updateDTO := ports.UpdatePublicationInput{Title: &newTitle}
	updatedPub, err := pubService.UpdatePublication(context.Background(), pub.ID, updateDTO)
	assert.NoError(t, err)
	assert.NotNil(t, updatedPub)
	assert.Equal(t, "Updated Awesome Title", updatedPub.Title)
	assert.Equal(t, "updated-awesome-title", updatedPub.Slug)
}
func TestPublicationService_Integration_Publishing(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	pubService := services.NewPublicationService(testutil.TestDB)
	author := models.User{FirstName: "Author", LoginEmail: "author@pub.com", Active: true}
	testutil.TestDB.Create(&author)
	pub := models.Publication{
		UserID:          author.ID,
		PublicationType: models.PublicationPost,
		Title:           "To Be Published",
		Slug:            "to-be-published",
		Content:         "Content.",
		Published:       false,
	}
	testutil.TestDB.Select("*").Create(&pub)
	assert.False(t, pub.Published)
	assert.Nil(t, pub.PublishedAt)
	published := true
	updateDTO := ports.UpdatePublicationInput{Published: &published}
	publishedPub, err := pubService.UpdatePublication(context.Background(), pub.ID, updateDTO)
	assert.NoError(t, err)
	assert.True(t, publishedPub.Published)
	assert.NotNil(t, publishedPub.PublishedAt)
	assert.WithinDuration(t, time.Now(), *publishedPub.PublishedAt, 5*time.Second)
}
func TestPublicationService_Integration_DeletePublication(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	pubService := services.NewPublicationService(testutil.TestDB)
	author := models.User{FirstName: "Author", LoginEmail: "author@pub.com", Active: true}
	testutil.TestDB.Create(&author)
	pub := models.Publication{
		UserID:          author.ID,
		PublicationType: models.PublicationArticle,
		Title:           "To Delete",
		Slug:            "to-delete",
		Content:         "Content.",
	}
	testutil.TestDB.Create(&pub)
	err := pubService.DeletePublication(context.Background(), pub.ID)
	assert.NoError(t, err)
}
