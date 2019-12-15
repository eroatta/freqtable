package main

import (
	"fmt"
	"log"

	"github.com/eroatta/freqtable/adapter/wordcount"
	"github.com/eroatta/freqtable/adapter/wordcount/cloner"
	"github.com/eroatta/freqtable/adapter/wordcount/miner"
)

func main() {
	config := wordcount.ProcessorConfig{
		Cloner: cloner.New(),
		Miner:  miner.NewCount(),
	}

	url := "https://github.com/src-d/go-siva"
	processor := wordcount.NewProcessor(config)
	frequencies, err := processor.Extract(url) //Build(url, config)
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

/*
type BuilderConfig struct {
	cloner wordcount.Cloner
	miner  wordcount.Miner
}*/

/*
var ErrCloningRepository = errors.New("Error while reading/cloning remote repository")

func Build(url string, config BuilderConfig) (map[string]int, error) {
	// cloning step
	_, filesc, err := wordcount.Clone(url, config.cloner)
	if err != nil {
		// TODO: improve error logging
		log.Println(fmt.Sprintf("Error reading repository %s: %v", url, err))
		return nil, ErrCloningRepository
	}

	// parsing step
	parsedc := wordcount.Parse(filesc)
	files := wordcount.Merge(parsedc)

	// mining step
	miningResults := wordcount.Mine(files, config.miner)
	countResults := miningResults.(miner.Count)

	return countResults.Results().(map[string]int), nil
}*/
