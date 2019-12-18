package entity

// FrequencyTable represents a frequency table, indluding its unique identifier,
// the related values and the error if any.
type FrequencyTable struct {
	ID     string
	Values map[string]int
}
