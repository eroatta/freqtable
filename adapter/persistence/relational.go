package persistence

import (
	"context"
	"errors"
	"fmt"
	"log"

	"database/sql"

	"github.com/eroatta/freqtable/entity"
	"github.com/eroatta/freqtable/repository"
)

var (
	ErrNoResults     = errors.New("No results for the given query")
	ErrUnexpected    = errors.New("Unexpected error performing the current operation")
	ErrMissingFields = errors.New("Missing mandatory fields")
)

type relational struct {
	db *sql.DB
}

// NewRelational creates a new FrequencyTableRepository backed up by a Relational Database.
func NewRelational(conn *sql.DB) repository.FrequencyTableRepository {
	return &relational{
		db: conn,
	}
}

func (r *relational) Save(ctx context.Context, ft entity.FrequencyTable) (int64, error) {
	if ft.Values == nil {
		return 0, ErrMissingFields
	}

	tx, err := r.db.Begin()
	if err != nil {
		// TODO: change log level
		log.Println(fmt.Sprintf("Error beginning a transaction: %v", err))
		return 0, ErrUnexpected
	}

	ftStmt, err := tx.PrepareContext(ctx,
		"INSERT INTO frequency_table(name, date_created) VALUES($1, $2) RETURNING id")
	if err != nil {
		// TODO: change log level
		log.Println(fmt.Sprintf("Error preparing statement 'INSERT INTO': %v", err))
		return 0, ErrUnexpected
	}

	var id int64
	err = ftStmt.QueryRowContext(ctx, ft.Name, ft.DateCreated).Scan(&id)
	if err != nil {
		// TODO: change log level
		log.Println(fmt.Sprintf("Error executing statement 'INSERT INTO frequency_table': %v", err))
		return 0, ErrUnexpected
	}

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO frequency_table_item(frequency_table_id, word, times) VALUES ($1, $2, $3)")
	if err != nil {
		// TODO: change log level
		log.Println(fmt.Sprintf("Error preparing statement 'INSERT INTO': %v", err))
		return 0, ErrUnexpected
	}

	for word, times := range ft.Values {
		if _, err = stmt.ExecContext(ctx, id, word, times); err != nil {
			// TODO: change log level
			defer tx.Rollback()
			return 0, ErrUnexpected
		}
	}

	if err = stmt.Close(); err != nil {
		// TODO: change log level
		return 0, ErrUnexpected
	}

	if err = tx.Commit(); err != nil {
		// TODO: change log level
		defer tx.Rollback()
		return 0, ErrUnexpected
	}

	return id, nil
}

func (r *relational) Get(ctx context.Context, ID int64) (entity.FrequencyTable, error) {
	//query := "SELECT id, \"name\", date_created, last_updated FROM frequency_table WHERE id=$1"
	// TODO: use prepared statements...
	query := "SELECT id FROM frequency_table WHERE id=$1"
	var frequencyTable entity.FrequencyTable

	row := r.db.QueryRowContext(ctx, query, ID)
	switch err := row.Scan(&frequencyTable.ID); err {
	case sql.ErrNoRows:
		return entity.FrequencyTable{}, ErrNoResults
	case nil:
		// continue
	default:
		// TODO: change log level
		log.Println(fmt.Sprintf("Error executing query '%s': %v", query, err))
		return entity.FrequencyTable{}, ErrUnexpected
	}

	// TODO: add query for values...

	return frequencyTable, nil
}
