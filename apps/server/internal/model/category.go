package model

import "github.com/google/uuid"

type Category struct {
	Base

	UserID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_categories_user_name" json:"user_id"`
	Name   string    `gorm:"type:varchar(120);not null;uniqueIndex:idx_categories_user_name" json:"name"`
	Color  string    `gorm:"type:varchar(32);not null;default:'#64748b'" json:"color"`

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
