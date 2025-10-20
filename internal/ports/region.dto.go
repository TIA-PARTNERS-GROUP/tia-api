package ports

import "github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"

type RegionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func MapRegionToResponse(region *models.Region) RegionResponse {
	return RegionResponse{
		ID:   region.ID,
		Name: region.Name,
	}
}
