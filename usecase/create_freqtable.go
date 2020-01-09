package usecase

import (
	"context"
	"time"

	"github.com/eroatta/freqtable/entity"
	"github.com/eroatta/freqtable/repository"
)

// CreateFrequencyTableUsecase defines the contract for the use cases related to the
// creation of one or several frenquency tables.
type CreateFrequencyTableUsecase interface {
	// Create creates a single frequency table.
	Create(ctx context.Context, url string) (entity.FrequencyTable, error)
}

// NewCreateFrequencyTableUsecase initializes a new CreateFrequencyTableUsecase handler
// with the given repositories.
func NewCreateFrequencyTableUsecase(wcr repository.WordCountRepository, ftr repository.FrequencyTableRepository) createFrequencyTableUsecase {
	return createFrequencyTableUsecase{
		wcr: wcr,
		ftr: ftr,
	}
}

type createFrequencyTableUsecase struct {
	wcr repository.WordCountRepository
	ftr repository.FrequencyTableRepository
}

// Create creates a new entity.FrequencyTable from the given URL.
func (uc createFrequencyTableUsecase) Create(ctx context.Context, url string) (entity.FrequencyTable, error) {
	ft := entity.FrequencyTable{
		Name:        url,
		DateCreated: time.Now(),
	}

	values, err := uc.wcr.Extract(url)
	if err != nil {
		return entity.FrequencyTable{}, err
	}
	ft.Values = values

	id, err := uc.ftr.Save(ctx, ft)
	if err != nil {
		return entity.FrequencyTable{}, err
	}
	ft.ID = id

	return ft, nil
}
