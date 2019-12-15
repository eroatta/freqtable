package wordcount

import (
	"errors"
	"fmt"
	"go/ast"
	"log"

	"github.com/eroatta/freqtable/adapter/wordcount/miner"
)

type Processor struct {
	config ProcessorConfig
}

type ProcessorConfig struct {
	ClonerFunc Cloner
	MinerFunc  Miner
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
}

var ErrCloningRepository = errors.New("Error while reading/cloning remote repository")

func (p Processor) Extract(url string) (map[string]int, error) {
	// cloning step
	_, filesc, err := step.Clone(url, p.config.ClonerFunc)
	if err != nil {
		// TODO: improve error logging
		log.Println(fmt.Sprintf("Error reading repository %s: %v", url, err))
		return nil, ErrCloningRepository
	}

	// parsing step
	parsedc := step.Parse(filesc)
	files := step.Merge(parsedc)

	// mining step
	miningResults := step.Mine(files, p.config.MinerFunc)
	countResults := miningResults.(miner.Count)

	return countResults.Results().(map[string]int), nil
}
