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
	frequencies, err := processor.Extract(url)
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
