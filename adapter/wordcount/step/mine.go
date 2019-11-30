package step

import (
	"go/ast"

	"github.com/eroatta/freqtable/adapter/wordcount"
)

// Miner interface is used to define a custom miner.
type Miner interface {
	// Name provides the name of the miner.
	Name() string
	// Visit applies the mining logic while traversing the Abstract Syntax Tree.
	Visit(node ast.Node) ast.Visitor
}

// Mine traverses each Abstract Syntax Tree and applies every given miner to extract
// the required pre-processing information. It returns the miner after work is done.
func Mine(parsed []wordcount.File, miner Miner) Miner {
	for _, f := range parsed {
		if f.AST == nil {
			continue
		}

		ast.Walk(miner, f.AST)
	}

	return miner
}
