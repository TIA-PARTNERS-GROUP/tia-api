package ports
import (
	"time"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)
type CreateBusinessInput struct {
	OperatorUserID   uint                    `json:"operator_user_id" validate:"required"`
	Name             string                  `json:"name" validate:"required,min=2,max=100"`
	Tagline          *string                 `json:"tagline" validate:"omitempty,max=100"`
	Website          *string                 `json:"website" validate:"omitempty,url"`
	ContactName      *string                 `json:"contact_name" validate:"omitempty,max=60"`
	ContactPhoneNo   *string                 `json:"contact_phone_no" validate:"omitempty,max=20"`
	ContactEmail     *string                 `json:"contact_email" validate:"omitempty,email"`
	Description      *string                 `json:"description"`
	Address          *string                 `json:"address" validate:"omitempty,max=100"`
	City             *string                 `json:"city" validate:"omitempty,max=60"`
	State            *string                 `json:"state" validate:"omitempty,max=60"`
	Country          *string                 `json:"country" validate:"omitempty,max=60"`
	PostalCode       *string                 `json:"postal_code" validate:"omitempty,max=20"`
	Value            *float64                `json:"value"`
	BusinessType     models.BusinessType     `json:"business_type" validate:"required"`
	BusinessCategory models.BusinessCategory `json:"business_category" validate:"required"`
	BusinessPhase    models.BusinessPhase    `json:"business_phase" validate:"required"`
}
type UpdateBusinessInput struct {
	Name             *string                  `json:"name" validate:"omitempty,min=2,max=100"`
	Tagline          *string                  `json:"tagline" validate:"omitempty,max=100"`
	Website          *string                  `json:"website" validate:"omitempty,url"`
	ContactName      *string                  `json:"contact_name" validate:"omitempty,max=60"`
	ContactPhoneNo   *string                  `json:"contact_phone_no" validate:"omitempty,max=20"`
	ContactEmail     *string                  `json:"contact_email" validate:"omitempty,email"`
	Description      *string                  `json:"description"`
	Address          *string                  `json:"address" validate:"omitempty,max=100"`
	City             *string                  `json:"city" validate:"omitempty,max=60"`
	State            *string                  `json:"state" validate:"omitempty,max=60"`
	Country          *string                  `json:"country" validate:"omitempty,max=60"`
	PostalCode       *string                  `json:"postal_code" validate:"omitempty,max=20"`
	Value            *float64                 `json:"value"`
	BusinessType     *models.BusinessType     `json:"business_type"`
	BusinessCategory *models.BusinessCategory `json:"business_category"`
	BusinessPhase    *models.BusinessPhase    `json:"business_phase"`
	Active           *bool                    `json:"active"`
}
type BusinessesFilter struct {
	BusinessType     *models.BusinessType     `form:"business_type"`
	BusinessCategory *models.BusinessCategory `form:"business_category"`
	BusinessPhase    *models.BusinessPhase    `form:"business_phase"`
	Active           *bool                    `form:"active"`
	OperatorUserID   *uint                    `form:"operator_user_id"`
	Search           *string                  `form:"search"`
}
type BusinessResponse struct {
	ID               uint                    `json:"id"`
	OperatorUserID   uint                    `json:"operator_user_id"`
	Name             string                  `json:"name"`
	Tagline          *string                 `json:"tagline,omitempty"`
	Website          *string                 `json:"website,omitempty"`
	ContactName      *string                 `json:"contact_name,omitempty"`
	ContactPhoneNo   *string                 `json:"contact_phone_no,omitempty"`
	ContactEmail     *string                 `json:"contact_email,omitempty"`
	Description      *string                 `json:"description,omitempty"`
	Address          *string                 `json:"address,omitempty"`
	City             *string                 `json:"city,omitempty"`
	State            *string                 `json:"state,omitempty"`
	Country          *string                 `json:"country,omitempty"`
	PostalCode       *string                 `json:"postal_code,omitempty"`
	Value            *float64                `json:"value,omitempty"`
	BusinessType     models.BusinessType     `json:"business_type"`
	BusinessCategory models.BusinessCategory `json:"business_category"`
	BusinessPhase    models.BusinessPhase    `json:"business_phase"`
	Active           bool                    `json:"active"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
	OperatorUser     *UserResponse           `json:"operator_user,omitempty"`
}
type BusinessStatsResponse struct {
	TotalProjects     int64 `json:"total_projects"`
	TotalPublications int64 `json:"total_publications"`
	TotalTags         int64 `json:"total_tags"`
}
func MapBusinessToResponse(business *models.Business) BusinessResponse {
	resp := BusinessResponse{
		ID:               business.ID,
		OperatorUserID:   business.OperatorUserID,
		Name:             business.Name,
		Tagline:          business.Tagline,
		Website:          business.Website,
		ContactName:      business.ContactName,
		ContactPhoneNo:   business.ContactPhoneNo,
		ContactEmail:     business.ContactEmail,
		Description:      business.Description,
		Address:          business.Address,
		City:             business.City,
		State:            business.State,
		Country:          business.Country,
		PostalCode:       business.PostalCode,
		Value:            business.Value,
		BusinessType:     business.BusinessType,
		BusinessCategory: business.BusinessCategory,
		BusinessPhase:    business.BusinessPhase,
		Active:           business.Active,
		CreatedAt:        business.CreatedAt,
		UpdatedAt:        business.UpdatedAt,
	}
	if business.OperatorUser.ID != 0 {
		operatorResp := MapUserToResponse(&business.OperatorUser)
		resp.OperatorUser = &operatorResp
	}
	return resp
}
