package wordcount

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse_OnClosedChannel_ShouldSendNoElements(t *testing.T) {
	filesc := make(chan File)
	close(filesc)

	parsedc := parse(filesc)

	var parsedFiles int
	for range parsedc {
		parsedFiles++
	}

	assert.Equal(t, 0, parsedFiles)
}

func TestParse_OnFileWithError_ShouldSendFileWithErrorMessage(t *testing.T) {
	filesc := make(chan File)
	go func() {
		filesc <- File{
			Name: "failing.go",
			Raw:  []byte("packaaage failing"),
		}
		close(filesc)
	}()

	parsedc := parse(filesc)

	files := make([]File, 0)
	for file := range parsedc {
		files = append(files, file)
	}

	assert.Equal(t, 1, len(files))
	assert.EqualError(t, files[0].Error, "failing.go:1:1: expected 'package', found packaaage")
}

func TestParse_OnTwoFiles_ShouldSendTwoParsedFilesWithSameFileset(t *testing.T) {
	filesc := make(chan File)
	go func() {
		filesc <- File{
			Name: "main.go",
			Raw:  []byte("package main"),
		}

		filesc <- File{
			Name: "test.go",
			Raw:  []byte("package test"),
		}

		close(filesc)
	}()

	parsedc := parse(filesc)

	files := make(map[string]File)
	for file := range parsedc {
		files[file.Name] = file
	}

	assert.Equal(t, 2, len(files))

	mainF := files["main.go"]
	assert.NotNil(t, mainF.AST)
	assert.NotNil(t, mainF.FileSet)
	assert.NoError(t, mainF.Error)

	testF := files["test.go"]
	assert.NotNil(t, testF.AST)
	assert.NotNil(t, testF.FileSet)
	assert.NoError(t, testF.Error)

	assert.Equal(t, mainF.FileSet, testF.FileSet)
}

func TestMerge_OnClosedChannel_ShouldReturnEmptyArray(t *testing.T) {
	parsedc := make(chan File)
	close(parsedc)

	got := merge(parsedc)

	assert.Empty(t, got)
}

func TestMerge_OnTwoFiles_ShouldReturnTwoFiles(t *testing.T) {
	parsedc := make(chan File)
	go func() {
		parsedc <- File{}
		parsedc <- File{}
		close(parsedc)
	}()

	got := merge(parsedc)

	assert.Equal(t, 2, len(got))
}
