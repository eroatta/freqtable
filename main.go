package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

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
	//storage := persistence.NewInMemory()
	db, err := newPostgresDB("localhost", 5432, "postgres", "postgres", "freqtable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	storage := persistence.NewRelational(db)

	createFreqTableUC := usecase.NewCreateFrequencyTableUsecase(processor, storage)

	ctx := context.Background()
	ft, err := createFreqTableUC.Create(ctx, url)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(fmt.Sprintf("Frequency Table - ID: %d - Name: %s - # Values: %d", ft.ID, ft.Name, len(ft.Values)))
	for token, count := range ft.Values {
		if len(token) == 1 {
			continue
		}

		log.Println(fmt.Sprintf("Repository: %s - Word: %s - Count: %d", url, token, count))
	}
}

func newPostgresDB(host string, port int, user string, password string, dbname string) (*sql.DB, error) {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error opening connection - %s -: %v", connInfo, err))
		return nil, err
	}
	// defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalln(fmt.Sprintf("Error pinging remote server: %v", err))
		return nil, err
	}

	return db, nil
}
