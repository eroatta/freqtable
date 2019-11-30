package step

import (
	"go/parser"
	"go/token"

	"github.com/eroatta/freqtable/adapter/wordcount"
)

// Parse parses a file and creates an Abstract Syntax Tree (AST) representation.
// It handles and returns a channel of code.File elements.
func Parse(filesc <-chan wordcount.File) chan wordcount.File {
	fset := token.NewFileSet()

	parsedc := make(chan wordcount.File)
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

// Merge joins files when necessary.
func Merge(parsedc <-chan wordcount.File) []wordcount.File {
	files := make([]wordcount.File, 0)
	for file := range parsedc {
		files = append(files, file)
	}

	return files
}
