package mysql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL driver is obviously needed
)

// A Storer implements the storage.Interface using a MySQL database backend
type Storer struct {
	db *sql.DB

	gopherTable  string
	ratingsTable string
}

const createGophersQuery = "CREATE TABLE IF NOT EXISTS `%s` (`hash` BINARY(32) NOT NULL, `gopher` BLOB NOT NULL, PRIMARY KEY (`hash`));"
const createRatingsQuery = "CREATE TABLE IF NOT EXISTS `%s` (`hash` BINARY(32) NOT NULL, `rating` TINYINT UNSIGNED NOT NULL, INDEX `ratings` (`hash`) );"

// NewStorerWithCustomTables creates a new Gopher storer using a given MySQL DSN
func NewStorerWithCustomTables(ctx context.Context, dsn string, gopherTable, ratingsTable string) (*Storer, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// TODO: validate gopherTable and ratingsTable

	if _, err := tx.Exec(fmt.Sprintf(createGophersQuery, gopherTable)); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(fmt.Sprintf(createRatingsQuery, ratingsTable)); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Storer{db: db, gopherTable: gopherTable, ratingsTable: ratingsTable}, nil
}
