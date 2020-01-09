package persistence

import (
	"context"
	"errors"
	"math/rand"

	"github.com/eroatta/freqtable/entity"
	"github.com/eroatta/freqtable/repository"
)

type memory struct {
	elements map[int64]entity.FrequencyTable
}

func NewInMemory() repository.FrequencyTableRepository {
	return &memory{
		elements: make(map[int64]entity.FrequencyTable),
	}
}

func (m *memory) Save(ctx context.Context, ft entity.FrequencyTable) (int64, error) {
	r := rand.New(rand.NewSource(99))
	ft.ID = r.Int63()
	m.elements[ft.ID] = ft

	return ft.ID, nil
}

func (m *memory) Get(ctx context.Context, id int64) (entity.FrequencyTable, error) {
	ft, ok := m.elements[id]
	if !ok {
		return entity.FrequencyTable{}, errors.New("not found")
	}

	return ft, nil
}
