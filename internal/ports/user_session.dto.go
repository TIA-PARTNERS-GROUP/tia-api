package ports
import (
	"time"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)
type SessionStatsResponse struct {
	UserID          uint  `json:"user_id"`
	TotalSessions   int64 `json:"total_sessions"`
	ActiveSessions  int64 `json:"active_sessions"`
	ExpiredSessions int64 `json:"expired_sessions"`
	RevokedSessions int64 `json:"revoked_sessions"`
}
type UserSessionResponse struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	IPAddress *string    `json:"ip_address,omitempty"`
	UserAgent *string    `json:"user_agent,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	IsActive  bool       `json:"is_active"`
	
	User UserResponse `json:"user"`
}
func MapToUserSessionResponse(session *models.UserSession) UserSessionResponse {
	isActive := session.RevokedAt == nil && time.Now().Before(session.ExpiresAt)
	return UserSessionResponse{
		ID:        session.ID,
		UserID:    session.UserID,
		IPAddress: session.IPAddress,
		UserAgent: session.UserAgent,
		CreatedAt: session.CreatedAt,
		ExpiresAt: session.ExpiresAt,
		RevokedAt: session.RevokedAt,
		IsActive:  isActive,
		User:      MapUserToResponse(&session.User),
	}
}
func MapToUserSessionsResponse(sessions []models.UserSession) []UserSessionResponse {
	responses := make([]UserSessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = MapToUserSessionResponse(&session)
	}
	return responses
}
