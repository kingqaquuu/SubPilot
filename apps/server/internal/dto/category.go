package dto

import (
	"time"

	"github.com/google/uuid"
)

type CategoryRequest struct {
	Name  string `json:"name" binding:"required,max=120"`
	Color string `json:"color" binding:"omitempty,max=32"`
}

type CategoryResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
