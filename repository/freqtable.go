package repository

import (
	"context"

	"github.com/eroatta/freqtable/entity"
)

// FrequencyTableRepository represents a repository capable of storing a given model.FrequencyTable.
type FrequencyTableRepository interface {
	// Get retrieves a model.FrequencyTable through the ID.
	Get(ctx context.Context, ID int64) (entity.FrequencyTable, error)
	// Save saves a model.FrequencyTable on the underlaying datasource.
	Save(ctx context.Context, ft entity.FrequencyTable) (int64, error)
}
