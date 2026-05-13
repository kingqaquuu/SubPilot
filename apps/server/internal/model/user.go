package model

type User struct {
	Base

	Email        string `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Name         string `gorm:"type:varchar(120);not null" json:"name"`
}
