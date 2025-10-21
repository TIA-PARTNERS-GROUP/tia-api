package services
import (
	"context"
	"errors"
	"time"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)
type UserSessionService struct {
	db *gorm.DB
}
func NewUserSessionService(db *gorm.DB) *UserSessionService {
	return &UserSessionService{db: db}
}
func (s *UserSessionService) GetSessionByID(ctx context.Context, sessionID uint) (*models.UserSession, error) {
	var session models.UserSession
	err := s.db.WithContext(ctx).
		Preload("User").
		First(&session, sessionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrSessionNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &session, nil
}
func (s *UserSessionService) GetActiveSessions(ctx context.Context, userID uint) ([]models.UserSession, error) {
	var sessions []models.UserSession
	err := s.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, time.Now()).
		Order("created_at desc").
		Find(&sessions).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return sessions, nil
}
func (s *UserSessionService) RevokeSession(ctx context.Context, sessionID uint) error {
	result := s.db.WithContext(ctx).
		Model(&models.UserSession{}).
		Where("id = ? AND revoked_at IS NULL", sessionID).
		Update("revoked_at", time.Now())
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrSessionNotFound
	}
	return nil
}
func (s *UserSessionService) RevokeAllSessions(ctx context.Context, userID uint, excludeSessionID *uint) error {
	query := s.db.WithContext(ctx).
		Model(&models.UserSession{}).
		Where("user_id = ? AND revoked_at IS NULL", userID)
	if excludeSessionID != nil {
		query = query.Where("id != ?", *excludeSessionID)
	}
	result := query.Update("revoked_at", time.Now())
	if result.Error != nil {
		return ports.ErrDatabase
	}
	return nil
}
func (s *UserSessionService) UpdateSessionActivity(ctx context.Context, sessionID uint, newExpiry time.Time) error {
	result := s.db.WithContext(ctx).
		Model(&models.UserSession{}).
		Where("id = ? AND revoked_at IS NULL", sessionID).
		Update("expires_at", newExpiry)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrSessionNotFound
	}
	return nil
}
func (s *UserSessionService) GetSessionStats(ctx context.Context, userID uint) (*ports.SessionStatsResponse, error) {
	var stats ports.SessionStatsResponse
	if err := s.db.WithContext(ctx).Model(&models.UserSession{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalSessions).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	if err := s.db.WithContext(ctx).Model(&models.UserSession{}).
		Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, time.Now()).
		Count(&stats.ActiveSessions).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	if err := s.db.WithContext(ctx).Model(&models.UserSession{}).
		Where("user_id = ? AND expires_at < ?", userID, time.Now()).
		Count(&stats.ExpiredSessions).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	if err := s.db.WithContext(ctx).Model(&models.UserSession{}).
		Where("user_id = ? AND revoked_at IS NOT NULL", userID).
		Count(&stats.RevokedSessions).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	stats.UserID = userID
	return &stats, nil
}
func (s *UserSessionService) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	result := s.db.WithContext(ctx).
		Where("expires_at < ? OR revoked_at IS NOT NULL", time.Now()).
		Delete(&models.UserSession{})
	if result.Error != nil {
		return 0, ports.ErrDatabase
	}
	return result.RowsAffected, nil
}
