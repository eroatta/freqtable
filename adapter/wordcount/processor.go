package wordcount

import (
	"errors"
	"fmt"
	"go/ast"
	"log"
)

type Processor struct {
	config ProcessorConfig
}

func NewProcessor(config ProcessorConfig) Processor {
	return Processor{
		config: config,
	}
}

type ProcessorConfig struct {
	Cloner Cloner
	Miner  Miner
}

// Cloner interface is used to define a custom cloner.
type Cloner interface {
	// Clone accesses a repository and clones it.
	Clone(url string) (Repository, error)
	// Filenames retrieves the names of the existing files on a repository.
	Filenames() ([]string, error)
	// File provides the bytes representation of a given file.
	File(name string) ([]byte, error)
}

// Miner interface is used to define a custom miner.
type Miner interface {
	// Name provides the name of the miner.
	Name() string
	// Visit applies the mining logic while traversing the Abstract Syntax Tree.
	Visit(node ast.Node) ast.Visitor
	// Results provides the mining results.
	Results() map[string]int
}

var ErrCloningRepository = errors.New("Error while reading/cloning remote repository")

func (p Processor) Extract(url string) (map[string]int, error) {
	// cloning step
	_, filesc, err := clone(url, p.config.Cloner)
	if err != nil {
		// TODO: improve error logging
		log.Println(fmt.Sprintf("Error reading repository %s: %v", url, err))
		return nil, ErrCloningRepository
	}

	// parsing & mining steps
	parsedc := parse(filesc)
	files := merge(parsedc)
	miningResults := mine(files, p.config.Miner)

	return miningResults.Results(), nil
}
