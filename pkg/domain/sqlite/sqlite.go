package sqlite

import (
	"database/sql"
	"errors"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
	"time"
)

// SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
	lg logger.ILogger
}

func NewSnippet(db *sql.DB, lg logger.ILogger) *SnippetModel {
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

	stmt := "INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, ?, ?)"

	result, err := m.DB.Exec(stmt, title, content, createdAt, expiresAt)
	if err != nil {
		m.lg.Error(err.Error())
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		m.lg.Error(err.Error())
		return 0, err
	}

	return int(id), nil
}

// Get will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	// Write the SQL statement we want to execute. Again, I've split it over two
	// lines for readability.
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > current_time AND id = ?`
	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	row := m.DB.QueryRow(stmt, id)
	// Initialize a pointer to a new zeroed Snippet struct.
	s := &models.Snippet{}
	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function check for that
		// error specifically, and return our own domain.ErrNoRecord error
		// instead.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	// If everything went OK then return the Snippet object.
	return s, nil
}

// Latest will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
