package repository

// WordCountRepository represents a repository capable of extracting the dictionary
// words count from a source code repository.
type WordCountRepository interface {
	// Extract extracts a map of words and counts from a source code repository.
	Extract(url string) (map[string]int, error)
}
