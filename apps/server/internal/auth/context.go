package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const userIDContextKey = "user_id"

var ErrMissingUserID = errors.New("missing authenticated user id")

func SetUserID(c *gin.Context, userID uuid.UUID) {
	c.Set(userIDContextKey, userID)
}

func UserID(c *gin.Context) (uuid.UUID, error) {
	value, ok := c.Get(userIDContextKey)
	if !ok {
		return uuid.Nil, ErrMissingUserID
	}

	userID, ok := value.(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrMissingUserID
	}

	return userID, nil
}
