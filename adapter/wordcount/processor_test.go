package wordcount_test

import (
	"errors"
	"go/ast"
	"testing"

	"github.com/eroatta/freqtable/adapter/wordcount"
	"github.com/stretchr/testify/assert"
)

func TestNewProcessor_ShouldReturnNewProcessor(t *testing.T) {
	processor := wordcount.NewProcessor(wordcount.ProcessorConfig{})

	assert.NotNil(t, processor)
}

func TestExtract_OnProcessorWithFailingCloningStep_ShouldReturnError(t *testing.T) {
	cloner := testCloner{
		err: errors.New("HTTP 404 Not Found"),
	}

	config := wordcount.ProcessorConfig{
		Cloner: cloner,
		Miner:  nil,
	}
	processor := wordcount.NewProcessor(config)
	_, err := processor.Extract("https://github.com/eroatta/freqtable")

	assert.EqualError(t, err, wordcount.ErrCloningRepository.Error())
}

func TestExtract_OnProcessorWithFailingParsingStep_ShouldReturnError(t *testing.T) {
	cloner := testCloner{
		repository: wordcount.Repository{
			Name: "freqtable",
			URL:  "https://github.com/eroatta/freqtable",
		},
		filenames: []string{"main.go"},
		files: map[string][]byte{
			"main.go": []byte("paaaaaackage main"),
		},
	}

	config := wordcount.ProcessorConfig{
		Cloner: cloner,
		Miner:  nil,
	}
	processor := wordcount.NewProcessor(config)
	_, err := processor.Extract("https://github.com/eroatta/freqtable")

	assert.EqualError(t, err, wordcount.ErrParsingFile.Error())
}
func TestExtract_OnProcessorWithOneFailedFileParsingStep_ShouldReturnReturnValidResults(t *testing.T) {
	cloner := testCloner{
		repository: wordcount.Repository{
			Name: "freqtable",
			URL:  "https://github.com/eroatta/freqtable",
		},
		filenames: []string{"main.go", "test.go"},
		files: map[string][]byte{
			"main.go": []byte("packaaaage main"),
			"test.go": []byte("package main"),
		},
	}

	miner := testMiner{
		results: map[string]int{
			"main": 1,
		},
	}

	config := wordcount.ProcessorConfig{
		Cloner: cloner,
		Miner:  miner,
	}
	processor := wordcount.NewProcessor(config)
	results, err := processor.Extract("https://github.com/eroatta/freqtable")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 1, results["main"])
}

func TestExtract_OnProcessor_ShouldReturnValidResults(t *testing.T) {
	cloner := testCloner{
		repository: wordcount.Repository{
			Name: "freqtable",
			URL:  "https://github.com/eroatta/freqtable",
		},
		filenames: []string{"main.go"},
		files: map[string][]byte{
			"main.go": []byte("package main"),
		},
	}

	miner := testMiner{
		results: map[string]int{
			"main": 1,
		},
	}

	config := wordcount.ProcessorConfig{
		Cloner: cloner,
		Miner:  miner,
	}
	processor := wordcount.NewProcessor(config)
	results, err := processor.Extract("https://github.com/eroatta/freqtable")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 1, results["main"])
}

type testCloner struct {
	repository wordcount.Repository
	filenames  []string
	files      map[string][]byte
	err        error
}

func (t testCloner) Clone(url string) (wordcount.Repository, error) {
	if t.err != nil {
		return wordcount.Repository{}, t.err
	}

	return t.repository, nil
}

func (t testCloner) Filenames() ([]string, error) {
	return t.filenames, nil
}

func (t testCloner) File(name string) ([]byte, error) {
	return t.files[name], nil
}

type testMiner struct {
	results map[string]int
}

func (t testMiner) Name() string {
	return "testMiner"
}

func (t testMiner) Visit(node ast.Node) ast.Visitor {
	return t
}

func (t testMiner) Results() map[string]int {
	return t.results
}
