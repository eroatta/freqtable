package main

import (
	"database/sql"
	"fmt"

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
	config := wordcount.ProcessorConfig{
		Cloner: cloner.New(),
		Miner:  miner.NewCount(),
	}

	//url := "https://github.com/src-d/go-siva"
	processor := wordcount.NewProcessor(config)
	//storage := persistence.NewInMemory()
	db, err := newPostgresDB("localhost", 5432, "postgres", "postgres", "freqtable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	storage := persistence.NewRelational(db)

	createFreqTableUC := usecase.NewCreateFrequencyTableUsecase(processor, storage)

	r := rest.NewServer(createFreqTableUC)
	r.Run()
}

func newPostgresDB(host string, port int, user string, password string, dbname string) (*sql.DB, error) {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.WithError(err).Fatal(fmt.Sprintf("error opening connection - %s", connInfo))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.WithError(err).Fatal(fmt.Sprintf("error pinging remote server"))
		return nil, err
	}

	return db, nil
}
