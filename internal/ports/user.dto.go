package ports
import (
	"time"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)
type UserCreationSchema struct {
	FirstName      string `json:"first_name" validate:"required,min=2,max=60"`
	LastName       string `json:"last_name" validate:"omitempty,min=2,max=60"`
	LoginEmail     string `json:"login_email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8,max=72"`
	ContactEmail   string `json:"contact_email" validate:"omitempty,email"`
	ContactPhoneNo string `json:"contact_phone_no" validate:"omitempty,max=20"`
	AdkSessionID   string `json:"adk_session_id" validate:"omitempty,max=128"`
}
type UserUpdateSchema struct {
	FirstName      *string `json:"first_name" validate:"omitempty,min=2,max=60"`
	LastName       *string `json:"last_name" validate:"omitempty,min=2,max=60"`
	LoginEmail     *string `json:"login_email" validate:"omitempty,email"`
	Password       *string `json:"password" validate:"omitempty,min=8,max=72"`
	ContactEmail   *string `json:"contact_email" validate:"omitempty,email"`
	ContactPhoneNo *string `json:"contact_phone_no" validate:"omitempty,max=20"`
	AdkSessionID   *string `json:"adk_session_id" validate:"omitempty,max=128"`
	EmailVerified  *bool   `json:"email_verified"`
	Active         *bool   `json:"active"`
}
type UserResponse struct {
	ID             uint      `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       *string   `json:"last_name"`
	LoginEmail     string    `json:"login_email"`
	ContactEmail   *string   `json:"contact_email"`
	ContactPhoneNo *string   `json:"contact_phone_no"`
	EmailVerified  bool      `json:"email_verified"`
	Active         bool      `json:"active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
func MapUserToResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		LoginEmail:     user.LoginEmail,
		ContactEmail:   user.ContactEmail,
		ContactPhoneNo: user.ContactPhoneNo,
		EmailVerified:  user.EmailVerified,
		Active:         user.Active,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}
