package persister

import (
	"database/sql"
	"errors"
	"snippetbox/pkg/domain/adapters/persister/sqlite"
	"snippetbox/pkg/domain/ports"
	"snippetbox/pkg/logger"
	"snippetbox/storing/store"
)

const (
	ModelInstanceSqlite int = iota
	ModelInstancePostgrest
)

var ErrUnsupportedModel = errors.New("unsupported model")

/*func NewSnippetsFactory(storeInstance int, lg *logger.ILogger, db *sql.DB) (models.ISnippet, error) {

	switch storeInstance {
	case store.StorageInstanceSqlite:
		return sqlite.NewSnippetRepository(db, lg), nil
	case store.StorageInstancePostgres:
		return postgres.NewSnippet(db, lg), nil
	default:
		return nil, ErrUnsupportedModel
	}
}*/

func NewRepositoryFactory(storeInstance int, lg *logger.ILogger,
	db *sql.DB) (*ports.IUserRepository, *ports.ISnippetRepository, error) {
	switch storeInstance {
	case store.StorageInstanceSqlite:
		var ur ports.IUserRepository = sqlite.NewUserRepository(db, lg)
		var sr ports.ISnippetRepository = sqlite.NewSnippetRepository(db, lg)

		return &ur, &sr, nil
	case store.StorageInstancePostgres:
		return nil, nil, nil
	default:
		return nil, nil, ErrUnsupportedModel
	}

}
