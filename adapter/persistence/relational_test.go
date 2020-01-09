package persistence_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eroatta/freqtable/adapter/persistence"
	"github.com/eroatta/freqtable/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewRelational_ShouldReturnNewFrequencyTableRepository(t *testing.T) {
	ftr := persistence.NewRelational(nil)

	assert.NotNil(t, ftr)
}

func TestGet_OnRelationalWhenSQLError_ShouldReturnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Unexpected error mocking a database connection: %v", err))
	}
	defer db.Close()
	mock.ExpectQuery("SELECT id FROM frequency_table WHERE id=(.+)").
		WithArgs(1234567890).
		WillReturnError(errors.New("Connection refused"))

	ftr := persistence.NewRelational(db)
	ft, err := ftr.Get(context.TODO(), 1234567890)

	assert.Empty(t, ft)
	assert.Equal(t, persistence.ErrUnexpected, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_OnRelationalWhenNonExistingFrequencyTable_ShouldReturnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Unexpected error mocking a database connection: %v", err))
	}
	defer db.Close()
	rows := mock.NewRows([]string{"id"})
	mock.ExpectQuery("SELECT id FROM frequency_table WHERE id=(.+)").
		WithArgs(1234567890).
		WillReturnRows(rows)

	ftr := persistence.NewRelational(db)
	ft, err := ftr.Get(context.TODO(), 1234567890)

	assert.Empty(t, ft)
	assert.Equal(t, persistence.ErrNoResults, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_OnRelationalWhenExistingFrequencyTable_ShouldReturnElement(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Unexpected error mocking a database connection: %v", err))
	}
	defer db.Close()
	rows := mock.NewRows([]string{"id"}).AddRow(1234567890)
	mock.ExpectQuery("SELECT id FROM frequency_table WHERE id=(.+)").
		WithArgs(1234567890).
		WillReturnRows(rows)

	ftr := persistence.NewRelational(db)
	ft, err := ftr.Get(context.TODO(), 1234567890)

	assert.Equal(t, int64(1234567890), ft.ID)
	assert.Empty(t, ft.Values)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSave_OnRelationalWhenMissingMandatoryValues_ShouldReturnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Unexpected error mocking a database connection: %v", err))
	}
	defer db.Close()

	ft := entity.FrequencyTable{}

	ftr := persistence.NewRelational(db)
	id, err := ftr.Save(context.TODO(), ft)

	assert.Equal(t, int64(0), id)
	assert.Equal(t, persistence.ErrMissingFields, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSave_OnRelationalWhenSQLError_ShouldReturnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Unexpected error mocking a database connection: %v", err))
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO frequency_table(.+) VALUES(.+) RETURNING id")
	now := time.Now()
	mock.ExpectQuery("INSERT INTO frequency_table(.+) VALUES(.+) RETURNING id").
		WithArgs("testname", now).
		WillReturnError(errors.New("sql: unexisting table"))

	ftr := persistence.NewRelational(db)

	ft := entity.FrequencyTable{
		Name:        "testname",
		DateCreated: now,
		Values: map[string]int{
			"cars":  1,
			"house": 3,
		},
	}
	id, err := ftr.Save(context.TODO(), ft)

	assert.Equal(t, int64(0), id)
	assert.Equal(t, persistence.ErrUnexpected, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSave_OnRelationalWhenErrorInsertingItems_ShouldReturnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Unexpected error mocking a database connection: %v", err))
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO frequency_table(.+) VALUES(.+) RETURNING id")
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(int64(1234567890))
	mock.ExpectQuery("INSERT INTO frequency_table(.+) VALUES(.+) RETURNING id").
		WithArgs("testname", now).
		WillReturnRows(rows)

	mock.ExpectPrepare("INSERT INTO frequency_table_item(.+) VALUES(.+)")
	mock.ExpectExec("INSERT INTO frequency_table_item(.+) VALUES(.+)").
		WithArgs(1234567890, "cars", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO frequency_table_item(.+) VALUES(.+)").
		WithArgs(1234567890, "house", 3).
		WillReturnError(errors.New("sql: invalid value"))
	mock.ExpectRollback()

	ftr := persistence.NewRelational(db)

	ft := entity.FrequencyTable{
		Name:        "testname",
		DateCreated: now,
		Values: map[string]int{
			"cars":  1,
			"house": 3,
		},
	}
	id, err := ftr.Save(context.TODO(), ft)

	assert.Equal(t, int64(0), id)
	assert.Equal(t, persistence.ErrUnexpected, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSave_OnRelationalWhenValidFrequencyTable_ShouldReturnNoError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Unexpected error mocking a database connection: %v", err))
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO frequency_table(.+) VALUES(.+) RETURNING id")
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(int64(1234567890))
	mock.ExpectQuery("INSERT INTO frequency_table(.+) VALUES(.+) RETURNING id").
		WithArgs("testname", now).
		WillReturnRows(rows)

	mock.ExpectPrepare("INSERT INTO frequency_table_item(.+) VALUES(.+)")
	mock.ExpectExec("INSERT INTO frequency_table_item(.+) VALUES(.+)").
		WithArgs(1234567890, "cars", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO frequency_table_item(.+) VALUES(.+)").
		WithArgs(1234567890, "house", 3).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	ftr := persistence.NewRelational(db)

	ft := entity.FrequencyTable{
		Name:        "testname",
		DateCreated: now,
		Values: map[string]int{
			"cars":  1,
			"house": 3,
		},
	}
	id, err := ftr.Save(context.TODO(), ft)

	assert.Equal(t, int64(1234567890), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
