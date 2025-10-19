package ports

import (
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

type CreatePublicationInput struct {
	UserID          uint                   `json:"user_id" validate:"required"`
	BusinessID      *uint                  `json:"business_id"`
	PublicationType models.PublicationType `json:"publication_type" validate:"required"`
	Title           string                 `json:"title" validate:"required,min=2,max=255"`
	Content         string                 `json:"content" validate:"required"`
	Excerpt         *string                `json:"excerpt"`
	Thumbnail       *string                `json:"thumbnail" validate:"omitempty,url"`
	VideoURL        *string                `json:"video_url" validate:"omitempty,url"`
	Published       *bool                  `json:"published"`
}

type UpdatePublicationInput struct {
	Title     *string `json:"title" validate:"omitempty,min=2,max=255"`
	Content   *string `json:"content"`
	Excerpt   *string `json:"excerpt"`
	Thumbnail *string `json:"thumbnail" validate:"omitempty,url"`
	VideoURL  *string `json:"video_url" validate:"omitempty,url"`
	Published *bool   `json:"published"`
}

type PublicationResponse struct {
	ID              uint                   `json:"id"`
	Slug            string                 `json:"slug"`
	PublicationType models.PublicationType `json:"publication_type"`
	Title           string                 `json:"title"`
	Excerpt         *string                `json:"excerpt,omitempty"`
	Content         string                 `json:"content"`
	Thumbnail       *string                `json:"thumbnail,omitempty"`
	VideoURL        *string                `json:"video_url,omitempty"`
	Published       bool                   `json:"published"`
	PublishedAt     *time.Time             `json:"published_at,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Author          UserResponse           `json:"author"`
	Business        *BusinessResponse      `json:"business,omitempty"`
}

func MapPublicationToResponse(pub *models.Publication) PublicationResponse {
	resp := PublicationResponse{
		ID:              pub.ID,
		Slug:            pub.Slug,
		PublicationType: pub.PublicationType,
		Title:           pub.Title,
		Excerpt:         pub.Excerpt,
		Content:         pub.Content,
		Thumbnail:       pub.Thumbnail,
		VideoURL:        pub.VideoURL,
		Published:       pub.Published,
		PublishedAt:     pub.PublishedAt,
		CreatedAt:       pub.CreatedAt,
		UpdatedAt:       pub.UpdatedAt,
	}

	if pub.User.ID != 0 {
		resp.Author = MapUserToResponse(&pub.User)
	}
	if pub.Business != nil {
		bizResp := MapBusinessToResponse(pub.Business)
		resp.Business = &bizResp
	}
	return resp
}
