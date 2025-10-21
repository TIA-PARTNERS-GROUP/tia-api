package services
import (
	"context"
	"errors"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)
type L2EResponseService struct {
	db *gorm.DB
}
func NewL2EResponseService(db *gorm.DB) *L2EResponseService {
	return &L2EResponseService{db: db}
}
func (s *L2EResponseService) CreateL2EResponse(ctx context.Context, data ports.CreateL2EResponseInput) (*models.L2EResponse, error) {
	if err := s.db.WithContext(ctx).Select("id").First(&models.User{}, data.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrUserNotFound
		}
		return nil, ports.ErrDatabase
	}
	response := models.L2EResponse{
		UserID:   data.UserID,
		Response: data.Response,
	}
	if err := s.db.WithContext(ctx).Create(&response).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return &response, nil
}
func (s *L2EResponseService) GetL2EResponsesForUser(ctx context.Context, userID uint) ([]models.L2EResponse, error) {
	var responses []models.L2EResponse
	err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("date_added desc").
		Find(&responses).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return responses, nil
}
