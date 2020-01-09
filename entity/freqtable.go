package entity

import "time"

// FrequencyTable represents a frequency table, indluding its unique identifier,
// the related values and the error if any.
type FrequencyTable struct {
	ID          int64
	Name        string
	DateCreated time.Time
	LastUpdated time.Time
	Values      map[string]int
}
