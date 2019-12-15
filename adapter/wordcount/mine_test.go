package wordcount_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/freqtable/adapter/wordcount"
	"github.com/stretchr/testify/assert"
)

func TestMine_OnNoFiles_ShouldReturnMinersWithoutResults(t *testing.T) {
	processed := wordcount.Mine([]wordcount.File{}, &miner{name: "empty"})
	emptyMiner, ok := processed.(*miner)

	assert.True(t, ok)
	assert.NotNil(t, emptyMiner)
	assert.Equal(t, 0, emptyMiner.visits)
}

func TestMine_OnFileWithNilAST_ShouldReturnMinersWithoutResults(t *testing.T) {
	processed := wordcount.Mine([]wordcount.File{{Name: "main.go"}}, &miner{name: "empty"})
	emptyMiner, ok := processed.(*miner)

	assert.True(t, ok)
	assert.NotNil(t, emptyMiner)
	assert.Equal(t, 0, emptyMiner.visits)
}

func TestMine_OnTwoMiners_ShouldReturnResultsBothMiners(t *testing.T) {
	/* Created AST:
	    0  *ast.File {
	    1  .  Doc: nil
	    2  .  Package: 1:1
	    3  .  Name: *ast.Ident {
	    4  .  .  NamePos: 1:9
	    5  .  .  Name: "main"
	    6  .  .  Obj: nil
	    7  .  }
	    8  .  Decls: nil
	    9  .  Scope: *ast.Scope {
	   10  .  .  Outer: nil
	   11  .  .  Objects: map[string]*ast.Object (len = 0) {}
	   12  .  }
	   13  .  Imports: nil
	   14  .  Unresolved: nil
	   15  .  Comments: nil
	   16  }
	*/

	testFileset := token.NewFileSet()

	ast1, _ := parser.ParseFile(testFileset, "main.go", `package main`, parser.AllErrors)
	file1 := wordcount.File{
		Name:    "main.go",
		AST:     ast1,
		FileSet: testFileset,
	}

	ast2, _ := parser.ParseFile(testFileset, "test.go", `package test`, parser.AllErrors)
	file2 := wordcount.File{
		Name:    "test.go",
		AST:     ast2,
		FileSet: testFileset,
	}

	testMiner := &miner{name: "first"}
	processed := wordcount.Mine([]wordcount.File{file1, file2}, testMiner)
	firstMiner, ok := processed.(*miner)

	assert.True(t, ok)
	assert.NotNil(t, firstMiner)
	assert.Equal(t, 8, firstMiner.visits)
}

type miner struct {
	name   string
	visits int
}

func (m *miner) Name() string {
	return m.name
}

func (m *miner) Visit(n ast.Node) ast.Visitor {
	m.visits++
	return m
}

func (m *miner) Results() map[string]int {
	return nil
}
