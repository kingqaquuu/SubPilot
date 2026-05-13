package migration

import (
	"fmt"

	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto`).Error; err != nil {
		return fmt.Errorf("enable pgcrypto extension: %w", err)
	}

	if err := db.AutoMigrate(model.All()...); err != nil {
		return fmt.Errorf("auto migrate models: %w", err)
	}

	return nil
}
