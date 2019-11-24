package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/eroatta/freqtable/cloner"
	"github.com/eroatta/freqtable/miner"
	"github.com/eroatta/freqtable/step"
)

func main() {
	config := BuilderConfig{
		cloner: cloner.New(),
		miner:  miner.NewCount(),
	}

	url := "https://github.com/src-d/go-siva"
	frequencies, err := Build(url, config)
	if err != nil {
		log.Fatalln(err)
	}

	for token, count := range frequencies {
		if len(token) == 1 {
			continue
		}

		log.Println(fmt.Sprintf("Repository: %s - Word: %s - Count: %d", url, token, count))
	}
}

type BuilderConfig struct {
	cloner step.Cloner
	miner  step.Miner
}

var ErrCloningRepository = errors.New("Error while reading/cloning remote repository")

func Build(url string, config BuilderConfig) (map[string]int, error) {
	// cloning step
	_, filesc, err := step.Clone(url, config.cloner)
	if err != nil {
		// TODO: improve error logging
		log.Println(fmt.Sprintf("Error reading repository %s: %v", url, err))
		return nil, ErrCloningRepository
	}

	// parsing step
	parsedc := step.Parse(filesc)
	files := step.Merge(parsedc)

	// mining step
	miningResults := step.Mine(files, config.miner)
	countResults := miningResults.(miner.Count)

	return countResults.Results().(map[string]int), nil
}
