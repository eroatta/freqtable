package main

import (
	_ "github.com/lib/pq"

	"github.com/eroatta/freqtable/adapter/persistence"
	"github.com/eroatta/freqtable/adapter/rest"
	"github.com/eroatta/freqtable/adapter/wordcount"
	"github.com/eroatta/freqtable/adapter/wordcount/cloner"
	"github.com/eroatta/freqtable/adapter/wordcount/miner"
	"github.com/eroatta/freqtable/usecase"
	log "github.com/sirupsen/logrus"
)

func main() {
	// processor configuration
	config := wordcount.ProcessorConfig{
		Cloner: cloner.New(),
		Miner:  miner.NewCount(),
	}
	processor := wordcount.NewProcessor(config)

	// storage configuration
	conn, deferrable, err := persistence.NewConnection("localhost", 5432, "postgres", "postgres", "freqtable")
	if err != nil {
		log.Fatalln(err)
	}
	defer deferrable()
	storage := persistence.NewPostgreSQL(conn)

	// rules engine configuration
	createFreqTableUC := usecase.NewCreateFrequencyTableUsecase(processor, storage)

	// REST controller
	r := rest.NewServer(createFreqTableUC)
	r.Run()
}
