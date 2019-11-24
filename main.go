package main

import (
	"fmt"
	"log"

	"github.com/eroatta/freqtable/cloner"
	"github.com/eroatta/freqtable/miner"
	"github.com/eroatta/freqtable/step"
)

func main() {
	newGoodMain("https://github.com/src-d/go-siva")
}

func newGoodMain(url string) {
	// cloning step
	_, filesc, err := step.Clone(url, cloner.New())
	if err != nil {
		log.Fatalf("Error reading repository %s: %v", url, err)
	}

	// parsing step
	parsedc := step.Parse(filesc)
	files := step.Merge(parsedc)

	// mining step
	countMiner := miner.NewCount()

	miningResults := step.Mine(files, countMiner)

	countResults := miningResults[countMiner.Name()].(miner.Count)
	freq := countResults.Results().(map[string]int)
	for token, count := range freq {
		if len(token) == 1 {
			continue
		}
		
		log.Println(fmt.Sprintf("Repository: %s - Word: %s - Count: %d", url, token, count))
	}
}
