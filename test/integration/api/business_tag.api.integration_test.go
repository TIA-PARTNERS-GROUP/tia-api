package main
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings" 
	"testing"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants" 
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)
func TestBusinessTagAPI_Integration_Lifecycle(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()
	
	user, token := CreateTestUserAndLogin(t, router, "tagowner@api.com", "ValidPassword123!")
	biz := models.Business{Name: "Tagged Biz", OperatorUserID: user.ID, BusinessType: "Other", BusinessCategory: "Mixed", BusinessPhase: "Growth"}
	testutil.TestDB.Create(&biz)
	
	constBizBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.BusinessBase
	constTagBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.TagsBase
	var createdTag ports.BusinessTagResponse
	t.Run("Create Tag", func(t *testing.T) {
		createDTO := ports.CreateBusinessTagInput{
			TagType:     models.BusinessTagService,
			Description: "API Development",
		}
		body, _ := json.Marshal(createDTO)
		
		createTagPath := strings.Replace(constants.AppRoutes.BusinessTags, ":id", fmt.Sprintf("%d", biz.ID), 1)
		url := constBizBase + createTagPath
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		json.Unmarshal(w.Body.Bytes(), &createdTag)
		assert.Equal(t, "API Development", createdTag.Description)
		assert.NotZero(t, createdTag.ID)
	})
	t.Run("Get Tags", func(t *testing.T) {
		
		getTagPath := strings.Replace(constants.AppRoutes.BusinessTags, ":id", fmt.Sprintf("%d", biz.ID), 1)
		url := constBizBase + getTagPath
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var tagsResponse ports.BusinessTagsResponse
		json.Unmarshal(w.Body.Bytes(), &tagsResponse)
		if assert.Equal(t, 1, tagsResponse.Count) {
			assert.Equal(t, "API Development", tagsResponse.Tags[0].Description)
		}
	})
	t.Run("Delete Tag", func(t *testing.T) {
		
		url := fmt.Sprintf("%s/%d", constTagBase, createdTag.ID)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
	t.Run("Verify Deletion", func(t *testing.T) {
		
		getTagPath := strings.Replace(constants.AppRoutes.BusinessTags, ":id", fmt.Sprintf("%d", biz.ID), 1)
		url := constBizBase + getTagPath
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var tagsResponse ports.BusinessTagsResponse
		json.Unmarshal(w.Body.Bytes(), &tagsResponse)
		assert.Equal(t, 0, tagsResponse.Count)
	})
}
