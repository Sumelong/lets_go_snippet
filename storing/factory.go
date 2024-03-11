package storing

import (
	"database/sql"
	"errors"
	"snippetbox/pkg"
	"snippetbox/storing/migration"
	"snippetbox/storing/store"
)

const (
	StoreInstancePostgres int = iota
	StoreInstanceSqlite
)

var ErrUnsupportedStore = errors.New("unsupported store")

type Store struct {
	instance int
	sql      *sql.DB
}

func NewStoreFactory(app *pkg.App) *sql.DB {

	switch app.StoreInstance {
	case StoreInstancePostgres:
		app.Logging.Info("beginning connection on postgres data store ")
		return store.NewStorePostgres(app.Logging)
	case StoreInstanceSqlite:
		app.Logging.Info("beginning connection on sqlite data store ")
		return store.NewStoreSqlite(app.Logging)
	default:
		panic(ErrUnsupportedStore)
	}
}

func RunMigration(app *pkg.App) {

	switch app.StoreInstance {
	case StoreInstancePostgres:
		app.Logging.Info("running migration on postgres data store ")
		migration.NewPostgresMigration(app)
	case StoreInstanceSqlite:
		app.Logging.Info("running migration on sqlite  data store ")
		migration.NewSqliteMigration(app)
	default:
		panic(ErrUnsupportedStore)
	}
}
