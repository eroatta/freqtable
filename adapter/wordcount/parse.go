package wordcount

import (
	"go/parser"
	"go/token"
)

// parse parses a file and creates an Abstract Syntax Tree (AST) representation.
// It handles and returns a channel of code.File elements.
func parse(filesc <-chan File) chan File {
	fset := token.NewFileSet()

	parsedc := make(chan File)
	go func() {
		for file := range filesc {
			node, err := parser.ParseFile(fset, file.Name, file.Raw, parser.ParseComments)

			file.AST = node
			file.FileSet = fset
			file.Error = err
			parsedc <- file
		}

		close(parsedc)
	}()

	return parsedc
}

// merge joins files when necessary.
func merge(parsedc <-chan File) []File {
	files := make([]File, 0)
	for file := range parsedc {
		files = append(files, file)
	}

	return files
}
