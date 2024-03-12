package store

import (
	"database/sql"
	"errors"
	"snippetbox/pkg"
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

func NewStoreFactory(storeInstance int, lg pkg.Logger) *sql.DB {

	switch storeInstance {
	case StorageInstancePostgres:
		lg.Info("beginning connection on postgres data store ")
		return NewStorePostgres(lg)
	case StorageInstanceSqlite:
		lg.Info("beginning connection on sqlite data store ")
		return NewStoreSqlite(lg)
	default:
		panic(ErrUnsupportedStore)
	}
}

func RunMigration(storeInstance int, db *sql.DB, lg pkg.Logger) {

	switch storeInstance {
	case StorageInstancePostgres:
		lg.Info("running migration on postgres data store ")
		migration.NewPostgresMigration(db, lg)
	case StorageInstanceSqlite:
		lg.Info("running migration on sqlite  data store ")
		migration.NewSqliteMigration(db, lg)
	default:
		panic(ErrUnsupportedStore)
	}
}
