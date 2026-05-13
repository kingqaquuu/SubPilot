package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/auth"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/response"
)

func Auth(tokens *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Error(c, http.StatusUnauthorized, "missing_token", "missing bearer token")
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			response.Error(c, http.StatusUnauthorized, "invalid_token", "invalid bearer token")
			c.Abort()
			return
		}

		userID, err := tokens.Parse(strings.TrimSpace(parts[1]))
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid_token", "invalid bearer token")
			c.Abort()
			return
		}

		auth.SetUserID(c, userID)
		c.Next()
	}
}
