package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"

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
	if err := godotenv.Load(); err != nil {
		log.WithError(err).Fatal("Error while loading the env configuration.")
	}

	conn, deferrable, err := persistence.NewConnection(os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
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
