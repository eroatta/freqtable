package entity_test

import (
	"errors"
	"testing"

	"github.com/eroatta/freqtable/entity"
	"github.com/stretchr/testify/assert"
)

func TestStatus_OnFrequencyTable_ShouldReturnExpectedValue(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"no_error", nil, "valid"},
		{"with_error", errors.New("my test error"), "invalid"},
	}

	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			ft := entity.FrequencyTable{
				ID:    "test",
				Error: fixture.err,
			}

			assert.Equal(t, fixture.expected, ft.Status())
		})
	}
}
