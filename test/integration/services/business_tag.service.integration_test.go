package main
import (
	"context"
	"testing"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)
func TestBusinessTagService_Integration_CreateAndGet(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessTagService := services.NewBusinessTagService(testutil.TestDB)
	user := models.User{FirstName: "Tagger", LoginEmail: "tagger@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business := models.Business{
		Name:             "Tagged Business",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeTechnology,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business)
	createDTO := ports.CreateBusinessTagInput{
		BusinessID:  business.ID,
		TagType:     models.BusinessTagClient,
		Description: "Enterprise Client",
	}
	createdTag, err := businessTagService.CreateBusinessTag(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdTag)
	assert.Equal(t, business.ID, createdTag.BusinessID)
	assert.Equal(t, models.BusinessTagClient, createdTag.TagType)
	assert.Equal(t, "Enterprise Client", createdTag.Description)
	assert.NotNil(t, createdTag.Business)
	assert.Equal(t, business.Name, createdTag.Business.Name)
	fetchedTag, err := businessTagService.GetBusinessTag(context.Background(), createdTag.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedTag)
	assert.Equal(t, createdTag.TagType, fetchedTag.TagType)
	assert.Equal(t, createdTag.Description, fetchedTag.Description)
	assert.Equal(t, createdTag.CreatedAt, fetchedTag.CreatedAt)
}
func TestBusinessTagService_Integration_DuplicatePrevention(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessTagService := services.NewBusinessTagService(testutil.TestDB)
	user := models.User{FirstName: "Tagger", LoginEmail: "tagger2@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business := models.Business{
		Name:             "Tagged Business 2",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeConsulting,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business)
	createDTO := ports.CreateBusinessTagInput{
		BusinessID:  business.ID,
		TagType:     models.BusinessTagService,
		Description: "Consulting Services",
	}
	_, err := businessTagService.CreateBusinessTag(context.Background(), createDTO)
	assert.NoError(t, err)
	_, err = businessTagService.CreateBusinessTag(context.Background(), createDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrBusinessTagAlreadyExists, err)
}
func TestBusinessTagService_Integration_UpdateAndDelete(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessTagService := services.NewBusinessTagService(testutil.TestDB)
	user := models.User{FirstName: "Tagger", LoginEmail: "tagger3@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business := models.Business{
		Name:             "Tagged Business 3",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeRetail,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business)
	createDTO := ports.CreateBusinessTagInput{
		BusinessID:  business.ID,
		TagType:     models.BusinessTagSpecialty,
		Description: "E-commerce",
	}
	tag, err := businessTagService.CreateBusinessTag(context.Background(), createDTO)
	assert.NoError(t, err)
	newDescription := "E-commerce & Mobile Apps"
	specialtyType := models.BusinessTagSpecialty
	updateDTO := ports.UpdateBusinessTagInput{
		TagType:     &specialtyType,
		Description: &newDescription,
	}
	updatedTag, err := businessTagService.UpdateBusinessTag(context.Background(), tag.ID, updateDTO)
	assert.NoError(t, err)
	assert.NotNil(t, updatedTag)
	assert.Equal(t, models.BusinessTagSpecialty, updatedTag.TagType)
	assert.Equal(t, newDescription, updatedTag.Description)
	tags, err := businessTagService.GetBusinessTags(context.Background(), business.ID, nil)
	assert.NoError(t, err)
	assert.Len(t, tags, 1)
	assert.Equal(t, tag.ID, tags[0].ID)
	assert.Equal(t, newDescription, tags[0].Description)
	err = businessTagService.DeleteBusinessTag(context.Background(), tag.ID)
	assert.NoError(t, err)
	_, err = businessTagService.GetBusinessTag(context.Background(), tag.ID)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrBusinessTagNotFound, err)
	tags, err = businessTagService.GetBusinessTags(context.Background(), business.ID, nil)
	assert.NoError(t, err)
	assert.Len(t, tags, 0)
}
func TestBusinessTagService_Integration_GetTagsByType(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessTagService := services.NewBusinessTagService(testutil.TestDB)
	user := models.User{FirstName: "Tagger", LoginEmail: "tagger4@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business := models.Business{
		Name:             "Multi-Tag Business",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeTechnology,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business)
	_, err := businessTagService.CreateBusinessTag(context.Background(), ports.CreateBusinessTagInput{
		BusinessID:  business.ID,
		TagType:     models.BusinessTagClient,
		Description: "Fortune 500",
	})
	assert.NoError(t, err)
	_, err = businessTagService.CreateBusinessTag(context.Background(), ports.CreateBusinessTagInput{
		BusinessID:  business.ID,
		TagType:     models.BusinessTagService,
		Description: "Cloud Migration",
	})
	assert.NoError(t, err)
	_, err = businessTagService.CreateBusinessTag(context.Background(), ports.CreateBusinessTagInput{
		BusinessID:  business.ID,
		TagType:     models.BusinessTagService,
		Description: "DevOps",
	})
	assert.NoError(t, err)
	_, err = businessTagService.CreateBusinessTag(context.Background(), ports.CreateBusinessTagInput{
		BusinessID:  business.ID,
		TagType:     models.BusinessTagSpecialty,
		Description: "AI/ML",
	})
	assert.NoError(t, err)
	allTags, err := businessTagService.GetBusinessTags(context.Background(), business.ID, nil)
	assert.NoError(t, err)
	assert.Len(t, allTags, 4)
	clientTags, err := businessTagService.GetTagsByType(context.Background(), business.ID, models.BusinessTagClient)
	assert.NoError(t, err)
	assert.Len(t, clientTags, 1)
	assert.Equal(t, "Fortune 500", clientTags[0].Description)
	serviceTags, err := businessTagService.GetTagsByType(context.Background(), business.ID, models.BusinessTagService)
	assert.NoError(t, err)
	assert.Len(t, serviceTags, 2)
	specialtyTags, err := businessTagService.GetTagsByType(context.Background(), business.ID, models.BusinessTagSpecialty)
	assert.NoError(t, err)
	assert.Len(t, specialtyTags, 1)
	assert.Equal(t, "AI/ML", specialtyTags[0].Description)
}
func TestBusinessTagService_Integration_SearchTags(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessTagService := services.NewBusinessTagService(testutil.TestDB)
	user := models.User{FirstName: "Tagger", LoginEmail: "tagger5@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business := models.Business{
		Name:             "Searchable Business",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeConsulting,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business)
	_, err := businessTagService.CreateBusinessTag(context.Background(), ports.CreateBusinessTagInput{
		BusinessID:  business.ID,
		TagType:     models.BusinessTagService,
		Description: "Digital Transformation",
	})
	assert.NoError(t, err)
	_, err = businessTagService.CreateBusinessTag(context.Background(), ports.CreateBusinessTagInput{
		BusinessID:  business.ID,
		TagType:     models.BusinessTagService,
		Description: "Business Process Optimization",
	})
	assert.NoError(t, err)
	_, err = businessTagService.CreateBusinessTag(context.Background(), ports.CreateBusinessTagInput{
		BusinessID:  business.ID,
		TagType:     models.BusinessTagClient,
		Description: "Healthcare Industry",
	})
	assert.NoError(t, err)
	transformationTags, err := businessTagService.SearchBusinessTags(context.Background(), business.ID, "Transformation", nil)
	assert.NoError(t, err)
	assert.Len(t, transformationTags, 1)
	assert.Equal(t, "Digital Transformation", transformationTags[0].Description)
	businessTags, err := businessTagService.SearchBusinessTags(context.Background(), business.ID, "Business", nil)
	assert.NoError(t, err)
	assert.Len(t, businessTags, 1)
	assert.Equal(t, "Business Process Optimization", businessTags[0].Description)
	serviceType := models.BusinessTagService
	serviceBusinessTags, err := businessTagService.SearchBusinessTags(context.Background(), business.ID, "Business", &serviceType)
	assert.NoError(t, err)
	assert.Len(t, serviceBusinessTags, 1)
	assert.Equal(t, "Business Process Optimization", serviceBusinessTags[0].Description)
}
func TestBusinessTagService_Integration_GetBusinessesByTag(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessTagService := services.NewBusinessTagService(testutil.TestDB)
	user := models.User{FirstName: "Tagger", LoginEmail: "tagger6@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business1 := models.Business{
		Name:             "Business A",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeTechnology,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business1)
	business2 := models.Business{
		Name:             "Business B",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeConsulting,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business2)
	_, err := businessTagService.CreateBusinessTag(context.Background(), ports.CreateBusinessTagInput{
		BusinessID:  business1.ID,
		TagType:     models.BusinessTagService,
		Description: "Enterprise Solutions",
	})
	assert.NoError(t, err)
	_, err = businessTagService.CreateBusinessTag(context.Background(), ports.CreateBusinessTagInput{
		BusinessID:  business2.ID,
		TagType:     models.BusinessTagService,
		Description: "Enterprise Solutions",
	})
	assert.NoError(t, err)
	businessesWithTag, err := businessTagService.GetBusinessesByTag(context.Background(), models.BusinessTagService, "Enterprise Solutions")
	assert.NoError(t, err)
	assert.Len(t, businessesWithTag, 2)
}
func TestBusinessTagService_Integration_Validation(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessTagService := services.NewBusinessTagService(testutil.TestDB)
	createDTO := ports.CreateBusinessTagInput{
		BusinessID:  999,
		TagType:     models.BusinessTagClient,
		Description: "Test Client",
	}
	_, err := businessTagService.CreateBusinessTag(context.Background(), createDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrBusinessNotFound, err)
}
func TestBusinessTagService_Integration_NonExistentTag(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessTagService := services.NewBusinessTagService(testutil.TestDB)
	_, err := businessTagService.GetBusinessTag(context.Background(), 999)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrBusinessTagNotFound, err)
	clientType := models.BusinessTagClient
	updateDTO := ports.UpdateBusinessTagInput{
		TagType: &clientType,
	}
	_, err = businessTagService.UpdateBusinessTag(context.Background(), 999, updateDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrBusinessTagNotFound, err)
	err = businessTagService.DeleteBusinessTag(context.Background(), 999)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrBusinessTagNotFound, err)
}
func TestBusinessTagService_Integration_UpdateDuplicatePrevention(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessTagService := services.NewBusinessTagService(testutil.TestDB)
	user := models.User{FirstName: "Tagger", LoginEmail: "tagger7@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business := models.Business{
		Name:             "Duplicate Test Business",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeTechnology,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business)
	_, err := businessTagService.CreateBusinessTag(context.Background(), ports.CreateBusinessTagInput{
		BusinessID:  business.ID,
		TagType:     models.BusinessTagService,
		Description: "Service A",
	})
	assert.NoError(t, err)
	tag2, err := businessTagService.CreateBusinessTag(context.Background(), ports.CreateBusinessTagInput{
		BusinessID:  business.ID,
		TagType:     models.BusinessTagService,
		Description: "Service B",
	})
	assert.NoError(t, err)
	serviceType := models.BusinessTagService
	duplicateDescription := "Service A"
	updateDTO := ports.UpdateBusinessTagInput{
		TagType:     &serviceType,
		Description: &duplicateDescription,
	}
	_, err = businessTagService.UpdateBusinessTag(context.Background(), tag2.ID, updateDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrBusinessTagAlreadyExists, err)
}
func TestBusinessTagService_Integration_UpdateNoData(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessTagService := services.NewBusinessTagService(testutil.TestDB)
	user := models.User{FirstName: "Tagger", LoginEmail: "tagger8@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business := models.Business{
		Name:             "No Update Business",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeTechnology,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business)
	createDTO := ports.CreateBusinessTagInput{
		BusinessID:  business.ID,
		TagType:     models.BusinessTagClient,
		Description: "Test Client",
	}
	tag, err := businessTagService.CreateBusinessTag(context.Background(), createDTO)
	assert.NoError(t, err)
	updateDTO := ports.UpdateBusinessTagInput{}
	_, err = businessTagService.UpdateBusinessTag(context.Background(), tag.ID, updateDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrNoUpdateData, err)
}
