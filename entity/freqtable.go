package entity

// FrequencyTable represents a frequency table, indluding its unique identifier,
// the related values and the error if any.
type FrequencyTable struct {
	ID     string
	Error  error
	Values map[string]int
}

// Status indicates if the FrequencyTable is valid or not.
func (ft FrequencyTable) Status() string {
	if ft.Error != nil {
		return "invalid"
	}

	return "valid"
}
