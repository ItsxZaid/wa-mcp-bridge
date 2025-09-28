package store

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type sqliteStore struct {
	db *sql.DB
	metaData MetadataStore
}

type Store interface {
	Metadata() MetadataStore
	Close() error
}

func New() (Store, error) {
	db, err := sql.Open("sqlite3", "./wa-mcp-bridge.db")
	if  err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := setUpMigration(db); err != nil {
		return nil, err
	}

	meta := &sqliteMetadataStore{db: db}

	return &sqliteStore{db: db, metaData: meta}, nil
}

func setUpMigration(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
        "file://db/migrations", 
        "sqlite3", 
        driver,
    )

	if err != nil {	
		return err
	}
	
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (s *sqliteStore) Metadata() MetadataStore {
	return s.metaData
}

func (s *sqliteStore) Close() error {
	return s.db.Close()
}