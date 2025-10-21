package middleware
import (
	"net/http"
	"strings"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"     
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services" 
	"github.com/gin-gonic/gin"
)
func AuthMiddleware(authService *services.AuthService, routes *constants.Routes) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}
		token := parts[1]
		user, session, err := authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}
		c.Set(routes.ContextKeyUser, user)
		c.Set(routes.ContextKeyUserID, user.ID)
		c.Set(routes.ContextKeySessionID, session.ID)
		c.Next()
	}
}
