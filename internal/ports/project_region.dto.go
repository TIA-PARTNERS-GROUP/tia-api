package ports

import "github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"

type AddProjectRegionInput struct {
	ProjectID uint   `json:"project_id" validate:"required"`
	RegionID  string `json:"region_id" validate:"required"`
}

type ProjectRegionResponse struct {
	ProjectID uint           `json:"project_id"`
	RegionID  string         `json:"region_id"`
	Region    RegionResponse `json:"region"`
}

func MapProjectRegionToResponse(pr *models.ProjectRegion) ProjectRegionResponse {
	resp := ProjectRegionResponse{
		ProjectID: pr.ProjectID,
		RegionID:  pr.RegionID,
	}
	if pr.Region.ID != "" {
		resp.Region = MapRegionToResponse(&pr.Region)
	}
	return resp
}
