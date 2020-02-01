package wordcount

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrCloningRepository indicates that an error occurred while retrieving the remote Github.com repository.
	ErrCloningRepository = errors.New("Error while reading/cloning remote repository")
	// ErrParsingFile indicates an error while converting the source code to its Abstract Syntax Tree representation.
	ErrParsingFile = errors.New("Error while parsing source code to AST")
)

// Processor handles the logic to extract the word count from a remote source code repository.
type Processor struct {
	config ProcessorConfig
}

// NewProcessor creates a new Processor based on the provided configuration.
func NewProcessor(config ProcessorConfig) Processor {
	return Processor{
		config: config,
	}
}

// Extract explores the source code and applies the processor-defined miner.
func (p Processor) Extract(url string) (map[string]int, error) {
	// cloning step
	_, filesc, err := clone(url, p.config.Cloner)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error reading repository %s", url))
		return nil, ErrCloningRepository
	}

	// parsing & mining steps
	parsedc := parse(filesc)
	files := merge(parsedc)
	for _, file := range files {
		if file.Error != nil {
			log.WithError(file.Error).Error(fmt.Sprintf("error when trying to parse file %s", file.Name))
			return nil, ErrParsingFile
		}
	}
	miningResults := mine(files, p.config.Miner)

	return miningResults.Results(), nil
}
