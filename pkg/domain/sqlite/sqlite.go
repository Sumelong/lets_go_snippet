package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
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
func (m *SnippetModel) Insert(title, content, expire string) (int, error) {

	exp := fmt.Sprintf("+%s days", expire)
	stmt := `INSERT INTO snippets (title, content, created, expires) 
			 VALUES(?, ?, CURRENT_TIMESTAMP, DATETIME(CURRENT_TIMESTAMP,?))`

	result, err := m.DB.Exec(stmt, title, content, exp)
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
	WHERE   id = ?`
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

// Remove take away will return a specific snippet based on its id.
func (m *SnippetModel) Remove(id int) (int, error) {
	// Write the SQL statement we want to execute. Again, I've split it over two
	// lines for readability.

	stmt := `DELETE  FROM snippets WHERE   id = ?`
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

	//where := time.Now()

	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > current_timestamp ORDER BY created DESC LIMIT 10`

	/*stmt := `SELECT id, title, content, created, expires FROM snippets
	ORDER BY created DESC LIMIT 10`
	*/
	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultset containing the result of
	// our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Initialize an empty slice to hold the models.Snippets objects.
	var snippets []*models.Snippet
	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.

	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &models.Snippet{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}
	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the Snippets slice.
	return snippets, nil
}
