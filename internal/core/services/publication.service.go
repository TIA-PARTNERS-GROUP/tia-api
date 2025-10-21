package services

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)

type PublicationService struct {
	db *gorm.DB
}

func NewPublicationService(db *gorm.DB) *PublicationService {
	return &PublicationService{db: db}
}
func generateSlug(title string) string {
	slug := strings.ToLower(title)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
func (s *PublicationService) CreatePublication(ctx context.Context, data ports.CreatePublicationInput) (*models.Publication, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, data.UserID).Error; err != nil {
		return nil, ports.ErrPublicationAuthorNotFound
	}
	slug := generateSlug(data.Title)
	published := false
	var publishedAt *time.Time
	if data.Published != nil && *data.Published {
		published = true
		now := time.Now()
		publishedAt = &now
	}
	publication := models.Publication{
		UserID:          data.UserID,
		BusinessID:      data.BusinessID,
		PublicationType: data.PublicationType,
		Title:           data.Title,
		Slug:            slug,
		Content:         data.Content,
		Excerpt:         data.Excerpt,
		Thumbnail:       data.Thumbnail,
		VideoURL:        data.VideoURL,
		Published:       published,
		PublishedAt:     publishedAt,
	}
	if err := s.db.WithContext(ctx).Create(&publication).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ports.ErrPublicationSlugExists
		}
		return nil, ports.ErrDatabase
	}
	return s.GetPublicationByID(ctx, publication.ID)
}
func (s *PublicationService) GetPublicationByID(ctx context.Context, id uint) (*models.Publication, error) {
	var pub models.Publication
	err := s.db.WithContext(ctx).Preload("User").Preload("Business").First(&pub, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrPublicationNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &pub, nil
}
func (s *PublicationService) GetPublicationBySlug(ctx context.Context, slug string) (*models.Publication, error) {
	var pub models.Publication
	err := s.db.WithContext(ctx).Preload("User").Preload("Business").Where("slug = ?", slug).First(&pub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrPublicationNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &pub, nil
}
func (s *PublicationService) UpdatePublication(ctx context.Context, id uint, data ports.UpdatePublicationInput) (*models.Publication, error) {
	_, err := s.GetPublicationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	updateData := make(map[string]interface{})
	if data.Title != nil {
		updateData["title"] = *data.Title
		updateData["slug"] = generateSlug(*data.Title)
	}
	if data.Content != nil {
		updateData["content"] = *data.Content
	}
	if data.Excerpt != nil {
		updateData["excerpt"] = *data.Excerpt
	}
	if data.Thumbnail != nil {
		updateData["thumbnail"] = *data.Thumbnail
	}
	if data.VideoURL != nil {
		updateData["video_url"] = *data.VideoURL
	}
	if data.Published != nil {
		updateData["published"] = *data.Published
		if *data.Published {
			updateData["published_at"] = time.Now()
		} else {
			updateData["published_at"] = nil
		}
	}
	if len(updateData) == 0 {
		return nil, ports.ErrNoUpdateData
	}
	if err := s.db.WithContext(ctx).Model(&models.Publication{ID: id}).Updates(updateData).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ports.ErrPublicationSlugExists
		}
		return nil, ports.ErrDatabase
	}
	return s.GetPublicationByID(ctx, id)
}
func (s *PublicationService) DeletePublication(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Publication{}, id)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrPublicationNotFound
	}
	return nil
}

func (s *PublicationService) FindAllPublications(ctx context.Context) ([]models.Publication, error) {
	var publications []models.Publication
	err := s.db.WithContext(ctx).
		Preload("User").
		Preload("Business").
		Order("published_at desc").
		Find(&publications).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return publications, nil
}
