package persistence

import "github.com/eroatta/freqtable/entity"

import "github.com/eroatta/freqtable/repository"

import "context"

import "errors"

type memory struct {
	elements map[string]entity.FrequencyTable
}

func NewInMemory() repository.FrequencyTableRepository {
	return &memory{
		elements: make(map[string]entity.FrequencyTable),
	}
}

func (m *memory) Save(ctx context.Context, ft entity.FrequencyTable) error {
	m.elements[ft.ID] = ft
	return nil
}

func (m *memory) Get(ctx context.Context, ID string) (entity.FrequencyTable, error) {
	ft, ok := m.elements[ID]
	if !ok {
		return entity.FrequencyTable{}, errors.New("not found")
	}

	return ft, nil
}
