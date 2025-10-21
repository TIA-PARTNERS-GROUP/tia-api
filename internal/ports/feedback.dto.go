package ports
import (
	"time"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)
type CreateFeedbackInput struct {
	Name    string `json:"name" validate:"required,min=2,max=120"`
	Email   string `json:"email" validate:"required,email"`
	Content string `json:"content" validate:"required"`
}
type FeedbackResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Content       string    `json:"content"`
	DateSubmitted time.Time `json:"date_submitted"`
}
func MapFeedbackToResponse(fb *models.Feedback) FeedbackResponse {
	return FeedbackResponse{
		ID:            fb.ID,
		Name:          fb.Name,
		Email:         fb.Email,
		Content:       fb.Content,
		DateSubmitted: fb.DateSubmitted,
	}
}
