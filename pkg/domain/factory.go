package domain

import (
	"database/sql"
	"errors"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/domain/postgres"
	"snippetbox/pkg/domain/sqlite"
	"snippetbox/pkg/logger"
	"snippetbox/storing/store"
)

const (
	ModelInstanceSqlite int = iota
	ModelInstancePostgrest
)

var ErrUnsupportedModel = errors.New("unsupported model")

func NewSnippetsFactory(storeInstance int, lg logger.ILogger, db *sql.DB) (models.ISnippet, error) {

	switch storeInstance {
	case store.StorageInstanceSqlite:
		return sqlite.NewSnippet(db, lg), nil
	case store.StorageInstancePostgres:
		return postgres.NewSnippet(db, lg), nil
	default:
		return nil, ErrUnsupportedModel
	}
}
