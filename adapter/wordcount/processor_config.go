package wordcount

import "go/ast"

// ProcessorConfig defines the properties available for configuration for a Processor.
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
