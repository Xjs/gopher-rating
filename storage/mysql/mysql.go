package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Xjs/gopher-rating/model"

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

// NewStorer creates a new Gopher storer using a given MySQL DSN
func NewStorer(ctx context.Context, dsn string) (*Storer, error) {
	return NewStorerWithCustomTables(ctx, dsn, "gophers", "ratings")
}

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

// Save implements the storage.Interface
func (s *Storer) Save(ctx context.Context, gopher *model.Gopher) error {
	if _, err := s.db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s VALUES (?, ?);", s.gopherTable), gopher.Hash, gopher.Raw); err != nil {
		return err
	}
	return nil
}

// Load implements the storage.Interface
func (s *Storer) Load(ctx context.Context, hash model.Hash) (*model.Gopher, error) {
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf("SELECT (`hash`, `gopher`) FROM %s WHERE hash = ?", s.gopherTable), hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gopher model.Gopher
	for rows.Next() {
		if err := rows.Scan(&gopher.Hash, &gopher.Raw); err != nil {
			return nil, err
		}
		return &gopher, nil
	}

	return nil, nil
}

// List implements the storage.Interface
func (s *Storer) List(ctx context.Context, start, count int) ([]model.Hash, error) {
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf("SELECT (`hash`) FROM %s ORDER BY `hash` LIMIT ?, ?", s.gopherTable), start, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hashes []model.Hash
	for rows.Next() {
		var h model.Hash
		if err := rows.Scan(&h); err != nil {
			return nil, err
		}
		hashes = append(hashes, h)
	}

	return hashes, nil
}

// Count implements the storage.Interface
func (s *Storer) Count(ctx context.Context) (int, error) {
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf("SELECT COUNT(`hash`) FROM %s", s.gopherTable))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
		return count, nil
	}

	return 0, errors.New("no count response from database")
}

// Rating implements the storage.Interface
func (s *Storer) Rating(ctx context.Context, hash model.Hash) (int, error) {
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf("SELECT AVG(`rating`) FROM %s WHERE hash = ?", s.ratingsTable), hash)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var rating int
	for rows.Next() {
		if err := rows.Scan(&rating); err != nil {
			return 0, err
		}
		return rating, nil
	}

	return 0, nil
}

// Rate implements the storage.Interface
func (s *Storer) Rate(ctx context.Context, hash model.Hash, rating int) error {
	if rating < 0 || rating > 5 {
		return errors.New("invalid rating")
	}
	if _, err := s.db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s VALUES (?, ?);", s.ratingsTable), hash, rating); err != nil {
		return err
	}
	return nil
}
