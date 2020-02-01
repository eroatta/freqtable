package persistence

import (
	"context"
	"errors"
	"fmt"

	"database/sql"

	"github.com/eroatta/freqtable/entity"
	"github.com/eroatta/freqtable/repository"

	log "github.com/sirupsen/logrus"
)

var (
	ErrNoResults     = errors.New("No results for the given query")
	ErrUnexpected    = errors.New("Unexpected error performing the current operation")
	ErrMissingFields = errors.New("Missing mandatory fields")
)

type postgresql struct {
	db *sql.DB
}

// NewPostgreSQL creates a new FrequencyTableRepository backed up by a Relational Database.
func NewPostgreSQL(conn *sql.DB) repository.FrequencyTableRepository {
	return &postgresql{
		db: conn,
	}
}

// NewConnection opens a connection between the application and the database and checks its validity.
// It returns the connection, a deferrable operation and error if present.
func NewConnection(host string, port int, user string, password string, dbname string) (*sql.DB, func(), error) {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.WithError(err).Fatal(fmt.Sprintf("error opening connection - %s", connInfo))
		return nil, func() {}, err
	}

	if err = db.Ping(); err != nil {
		log.WithError(err).Fatal(fmt.Sprintf("error pinging remote server"))
		return nil, func() {}, err
	}

	return db, func() { db.Close() }, nil
}

func (r *postgresql) Save(ctx context.Context, ft entity.FrequencyTable) (int64, error) {
	if ft.Values == nil {
		return 0, ErrMissingFields
	}

	tx, err := r.db.Begin()
	if err != nil {
		log.WithField("error", err).Error("error beginning a transaction")
		return 0, ErrUnexpected
	}

	ftStmt, err := tx.PrepareContext(ctx,
		"INSERT INTO frequency_table(name, date_created) VALUES($1, $2) RETURNING id")
	if err != nil {
		log.WithField("error", err).Error("error preparing statement for frequency_table insertion")
		return 0, ErrUnexpected
	}

	var id int64
	err = ftStmt.QueryRowContext(ctx, ft.Name, ft.DateCreated).Scan(&id)
	if err != nil {
		log.WithField("error", err).Error("error inserting new frequency_table record")
		return 0, ErrUnexpected
	}

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO frequency_table_item(frequency_table_id, word, times) VALUES ($1, $2, $3)")
	if err != nil {
		log.WithField("error", err).Error("error preparing statement for frequency_table_item insertion")
		return 0, ErrUnexpected
	}

	for word, times := range ft.Values {
		if _, err = stmt.ExecContext(ctx, id, word, times); err != nil {
			log.WithField("error", err).Error("error inserting new frequency_table_item record")
			defer tx.Rollback()
			return 0, ErrUnexpected
		}
	}

	if err = tx.Commit(); err != nil {
		log.WithField("error", err).Error("error commiting a transaction")
		defer tx.Rollback()
		return 0, ErrUnexpected
	}

	return id, nil
}

func (r *postgresql) Get(ctx context.Context, ID int64) (entity.FrequencyTable, error) {
	query := "SELECT id, \"name\", date_created, last_updated FROM frequency_table WHERE id=$1"
	ftGetStmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.WithError(err).Error("error preparing frequency_table select statement")
		return entity.FrequencyTable{}, ErrUnexpected
	}

	var frequencyTable entity.FrequencyTable
	row := ftGetStmt.QueryRowContext(ctx, ID)
	switch err := row.Scan(&frequencyTable.ID,
		&frequencyTable.Name,
		&frequencyTable.DateCreated,
		&frequencyTable.LastUpdated); err {
	case sql.ErrNoRows:
		return entity.FrequencyTable{}, ErrNoResults
	case nil:
		// continue
	default:
		log.WithError(err).Error("error executing select on frequency_table")
		return entity.FrequencyTable{}, ErrUnexpected
	}

	itemsQuery := "SELECT word, times FROM frequency_table_item WHERE frequency_table_id=$1"
	itemsSelectStmt, err := r.db.PrepareContext(ctx, itemsQuery)
	if err != nil {
		log.WithError(err).Error("error preparing frequency_table_item select statement")
		return entity.FrequencyTable{}, ErrUnexpected
	}

	rows, err := itemsSelectStmt.QueryContext(ctx, frequencyTable.ID)
	if err != nil {
		log.WithError(err).Error("error executing select on frequency_table_item")
	}
	defer rows.Close()

	frequencyTable.Values = make(map[string]int)
	for rows.Next() {
		var word string
		var times int
		if err := rows.Scan(&word, &times); err != nil {
			log.WithError(err).Error("error scanning row results")
			return entity.FrequencyTable{}, ErrUnexpected
		}
		frequencyTable.Values[word] = times
	}

	return frequencyTable, nil
}
