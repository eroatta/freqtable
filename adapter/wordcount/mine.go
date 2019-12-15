package wordcount

import "go/ast"

// Mine traverses each Abstract Syntax Tree and applies every given miner to extract
// the required pre-processing information. It returns the miner after work is done.
func Mine(parsed []File, miner Miner) Miner {
	for _, f := range parsed {
		if f.AST == nil {
			continue
		}

		ast.Walk(miner, f.AST)
	}

	return miner
}
