package code

import (
	"go/ast"
	"go/token"
)

// Repository holds information of a GitHub repository.
type Repository struct {
	Name        string
	URL         string
	Hash        string
	DateCreated string
	Error       error
}

// File represents a file on a code.Repository, and contains a raw representation and
// a ast.File representation.
type File struct {
	Name    string
	Raw     []byte
	AST     *ast.File
	FileSet *token.FileSet
	Error   error
}
