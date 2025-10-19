package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/TIA-PARTNERS-GROUP/tia-api/pkg/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db                   *gorm.DB
	sessionCleanupTicker *time.Ticker
	quitChan             chan struct{}
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		db:       db,
		quitChan: make(chan struct{}),
	}
}

func (s *AuthService) Login(ctx context.Context, data ports.LoginInput, ipAddress, userAgent *string) (*ports.LoginResponse, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where("login_email = ?", data.LoginEmail).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrInvalidCredentials
		}
		return nil, ports.ErrDatabase
	}

	if !user.Active {
		return nil, ports.ErrAccountDeactivated
	}

	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return nil, ports.ErrInvalidCredentials
	}

	if err := utils.VerifyPassword(data.Password, *user.PasswordHash); err != nil {
		return nil, ports.ErrInvalidCredentials
	}

	session := models.UserSession{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // Set a preliminary expiry
	}
	if ipAddress != nil {
		session.IPAddress = ipAddress
	}
	if userAgent != nil {
		session.UserAgent = userAgent
	}
	if err := s.db.WithContext(ctx).Create(&session).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	token, expiry, err := utils.GenerateToken(user.ID, session.ID, user.LoginEmail)
	if err != nil {
		return nil, ports.ErrTokenGeneration
	}

	sessionUpdate := map[string]interface{}{
		"token_hash": utils.HashToken(token),
		"expires_at": expiry,
	}
	if err := s.db.WithContext(ctx).Model(&session).Updates(sessionUpdate).Error; err != nil {
		s.db.Delete(&session)
		return nil, ports.ErrDatabase
	}

	response := &ports.LoginResponse{
		User:      ports.MapUserToResponse(&user),
		Token:     token,
		SessionID: session.ID,
		ExpiresAt: expiry,
		TokenType: "Bearer",
	}

	return response, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID, userID uint) (bool, error) {
	result := s.db.WithContext(ctx).
		Model(&models.UserSession{}).
		Where("id = ? AND user_id = ? AND revoked_at IS NULL", sessionID, userID).
		Update("revoked_at", time.Now())

	if result.Error != nil {
		return false, ports.ErrDatabase
	}

	return result.RowsAffected > 0, nil
}

func (s *AuthService) LogoutAll(ctx context.Context, userID uint) (int64, error) {
	result := s.db.WithContext(ctx).
		Model(&models.UserSession{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", time.Now())

	if result.Error != nil {
		return 0, ports.ErrDatabase
	}
	return result.RowsAffected, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*models.User, *models.UserSession, error) {
	claims, err := utils.VerifyToken(token)
	if err != nil {
		return nil, nil, ports.ErrInvalidToken
	}

	var session models.UserSession
	err = s.db.WithContext(ctx).
		Preload("User").
		Where("id = ? AND user_id = ? AND revoked_at IS NULL", claims.SessionID, claims.UserID).
		First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ports.ErrInvalidSession
		}
		return nil, nil, ports.ErrDatabase
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, nil, ports.ErrInvalidSession
	}

	if !session.User.Active {
		return nil, nil, ports.ErrAccountDeactivated
	}

	return &session.User, &session, nil
}

func (s *AuthService) GetUserSessions(ctx context.Context, userID uint) ([]models.UserSession, error) {
	var sessions []models.UserSession
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, time.Now()).
		Order("created_at desc").
		Find(&sessions).Error

	if err != nil {
		return nil, ports.ErrDatabase
	}
	return sessions, nil
}

func (s *AuthService) StartSessionCleanup() {
	log.Println("Initializing session cleanup job...")
	go s.cleanupExpiredSessions()

	s.sessionCleanupTicker = time.NewTicker(24 * time.Hour)
	go func() {
		for {
			select {
			case <-s.sessionCleanupTicker.C:
				log.Println("Running scheduled session cleanup job...")
				s.cleanupExpiredSessions()
			case <-s.quitChan:
				s.sessionCleanupTicker.Stop()
				return
			}
		}
	}()
	log.Println("Session cleanup job started.")
}

func (s *AuthService) StopSessionCleanup() {
	close(s.quitChan)
	log.Println("Session cleanup job stopped.")
}

func (s *AuthService) cleanupExpiredSessions() {
	result := s.db.
		Where("expires_at < ? OR revoked_at IS NOT NULL", time.Now()).
		Delete(&models.UserSession{})

	if result.Error != nil {
		log.Printf("Session cleanup failed: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Session cleanup completed: %d sessions removed", result.RowsAffected)
	}
}
