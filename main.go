package main

import (
	"context"
	"fmt"
	"log"

	"github.com/eroatta/freqtable/adapter/persistence"
	"github.com/eroatta/freqtable/adapter/wordcount"
	"github.com/eroatta/freqtable/adapter/wordcount/cloner"
	"github.com/eroatta/freqtable/adapter/wordcount/miner"
	"github.com/eroatta/freqtable/usecase"
)

func main() {
	config := wordcount.ProcessorConfig{
		Cloner: cloner.New(),
		Miner:  miner.NewCount(),
	}

	url := "https://github.com/src-d/go-siva"
	processor := wordcount.NewProcessor(config)
	storage := persistence.NewInMemory()

	createFreqTableUC := usecase.NewCreateFrequencyTableUsecase(processor, storage)

	ctx := context.Background()
	ft, err := createFreqTableUC.Create(ctx, url)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(fmt.Sprintf("Frequency Table - ID: %s - # Values: %d", ft.ID, len(ft.Values)))
	for token, count := range ft.Values {
		if len(token) == 1 {
			continue
		}

		log.Println(fmt.Sprintf("Repository: %s - Word: %s - Count: %d", url, token, count))
	}
}
