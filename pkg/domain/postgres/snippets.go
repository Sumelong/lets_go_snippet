package postgres

import (
	"database/sql"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
	"time"
)

// SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
	lg logger.ILogger
}

func NewSnippet(db *sql.DB, lg *logger.Logger) *SnippetModel {
	return &SnippetModel{
		DB: db,
		lg: lg,
	}
}

// Insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {

	createdAt := time.Now()
	duration, err := time.ParseDuration(expires)
	expiresAt := time.Now().Add(time.Hour * 24 * duration)

	stmt := `INSERT INTO snippets (title, content, created, expires)
				VALUES($1, $2, $3, $4  )`

	result, err := m.DB.Exec(stmt, title, content, createdAt, expiresAt)
	if err != nil {

		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// Remove take away will return a specific snippet based on its id.
func (m *SnippetModel) Remove(id int) (int, error) {
	// Write the SQL statement we want to execute. Again, I've split it over two
	// lines for readability.

	stmt := `DELETE  FROM snippets WHERE   id = $1`
	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	res, err := m.DB.Exec(stmt, id)
	row, _ := res.RowsAffected()
	if err != nil {
		return 0, models.ErrNoRecord

	}

	// If everything went OK then return the Snippet object.
	return int(row), nil
}

// Latest will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
