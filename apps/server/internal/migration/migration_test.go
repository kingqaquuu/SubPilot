package migration

import (
	"testing"

	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
)

func TestMigrationModelCoverage(t *testing.T) {
	if got := len(model.All()); got != 5 {
		t.Fatalf("expected migration coverage for 5 models, got %d", got)
	}
}
