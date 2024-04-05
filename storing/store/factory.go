package store

import (
	"database/sql"
	"errors"
	"snippetbox/pkg/logger"
	"snippetbox/storing/migration"
)

const (
	StorageInstancePostgres int = iota
	StorageInstanceSqlite
)

var ErrUnsupportedStore = errors.New("unsupported store")

type Store struct {
	instance int
	sql      *sql.DB
}

func NewStoreFactory(storeInstance int, lg *logger.ILogger) *sql.DB {

	switch storeInstance {
	case StorageInstancePostgres:
		return NewStorePostgres(lg)
	case StorageInstanceSqlite:
		return NewStoreSqlite(lg)
	default:
		panic(ErrUnsupportedStore)
	}
}

func RunMigration(storeInstance int, db *sql.DB, lg *logger.ILogger) {

	switch storeInstance {
	case StorageInstancePostgres:

		migration.NewPostgresMigration(db, lg)
	case StorageInstanceSqlite:

		migration.NewSqliteMigration(db, lg)
	default:
		panic(ErrUnsupportedStore)
	}
}
