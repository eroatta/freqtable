package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/eroatta/freqtable/entity"
	"github.com/eroatta/freqtable/usecase"
	"github.com/stretchr/testify/assert"
)

func TestNewCreateFrequencyTableUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := usecase.NewCreateFrequencyTableUsecase(nil, nil)

	assert.NotNil(t, uc)
}

func TestCreate_OnCreateFrequencyTableUsecase_ShouldCreateFrequencyTable(t *testing.T) {
	wcr := testWordCountRepository{
		extractions: map[string]map[string]int{
			"https://github.com/eroatta/freqtable": map[string]int{
				"frequency": 2,
				"table":     3,
			},
		},
		err: nil,
	}

	ftr := testFrequencyTableRepository{
		err: nil,
	}

	uc := usecase.NewCreateFrequencyTableUsecase(wcr, ftr)
	ft, err := uc.Create(context.TODO(), "https://github.com/eroatta/freqtable")

	assert.NotNil(t, ft)
	assert.NoError(t, err)
	assert.Equal(t, "60cc02e02b67414b493033995baf772f", ft.ID)
	assert.Equal(t, 2, len(ft.Values))
	assert.Equal(t, 2, ft.Values["frequency"])
	assert.Equal(t, 3, ft.Values["table"])
}

func TestCreate_OnCreateFrequencyTableUsecase_WhenErrorCounting_ShouldReturnError(t *testing.T) {
	wcr := testWordCountRepository{
		extractions: map[string]map[string]int{},
		err:         errors.New("error while extracting"),
	}

	uc := usecase.NewCreateFrequencyTableUsecase(wcr, nil)
	ft, err := uc.Create(context.TODO(), "https://github.com/eroatta/freqtable")

	assert.NotNil(t, ft)
	assert.EqualError(t, err, "error while extracting")
	assert.Equal(t, entity.FrequencyTable{}, ft)
}

func TestCreate_OnCreateFrequencyTableUsecase_WhenSavingResults_ShouldReturnError(t *testing.T) {
	wcr := testWordCountRepository{
		extractions: map[string]map[string]int{
			"https://github.com/eroatta/freqtable": map[string]int{
				"frequency": 2,
				"table":     3,
			},
		},
		err: nil,
	}

	ftr := testFrequencyTableRepository{
		err: errors.New("error while persisting"),
	}

	uc := usecase.NewCreateFrequencyTableUsecase(wcr, ftr)
	ft, err := uc.Create(context.TODO(), "https://github.com/eroatta/freqtable")

	assert.NotNil(t, ft)
	assert.EqualError(t, err, "error while persisting")
	assert.Equal(t, entity.FrequencyTable{}, ft)
}

type testWordCountRepository struct {
	extractions map[string]map[string]int
	err         error
}

func (twc testWordCountRepository) Extract(url string) (map[string]int, error) {
	if val, ok := twc.extractions[url]; ok {
		return val, nil
	}

	return nil, twc.err
}

type testFrequencyTableRepository struct {
	frequencyTable entity.FrequencyTable
	err            error
}

func (tft testFrequencyTableRepository) Get(ctx context.Context, ID string) (entity.FrequencyTable, error) {
	return tft.frequencyTable, tft.err
}

func (tft testFrequencyTableRepository) Save(ctx context.Context, ft entity.FrequencyTable) error {
	return tft.err
}
