package repository

import (
	"context"

	"github.com/eroatta/freqtable/entity"
)

// FrequencyTableRepository represents a repository capable of storing a given model.FrequencyTable.
type FrequencyTableRepository interface {
	// GetByID retrieves a model.FrequencyTable though the ID.
	Get(ctx context.Context, ID string) (entity.FrequencyTable, error)
	// Save saves a model.FrequencyTable on the underlaying datasource.
	Save(ctx context.Context, ft entity.FrequencyTable) error
}
