package migration

import (
	"database/sql"
	"snippetbox/pkg/logger"
)

func NewSqliteMigration(db *sql.DB, lg *logger.Logger) {

	defer db.Close()

	// Create your table(s) here. This is just an example.
	_, err := db.Exec(`
		CREATE TABLE snippets (
			id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
			title VARCHAR(100) NOT NULL,
			content TEXT NOT NULL,
			created DATETIME NOT NULL,
			expires DATETIME NOT NULL);`)
	if err != nil {
		lg.Error(err.Error())
		return
	}

}
