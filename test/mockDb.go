package test

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"snippetbox/pkg/logger"
	"snippetbox/storing/migration"
	"testing"
)

func MocKDB(t *testing.T, lg *logger.ILogger) (*sql.DB, func()) {

	// Create a temporary file for the SQLite database
	prefix := "testDb-*.db"

	tmpFile, err := os.CreateTemp(".", prefix)
	if err != nil {
		t.Errorf("Error creating temporary file: %v\n", err)
		return nil, nil
	}

	// Connect to the SQLite database
	db, err := sql.Open("sqlite", tmpFile.Name())
	if err != nil {
		t.Fatalf("Error connecting to the database: %v\n", err)

	}
	err = db.Ping()
	if err != nil {
		t.Fatalf("Error opening database: %v\n", err)
	}

	// Create tables and perform test operations
	migration.NewSqliteMigration(db, lg)

	// Clean up the temporary file when done
	cleanUp := func() {

		//drop table from data base
		_, err = db.Exec(` 
				DROP TABLE users;
				DROP TABLE snippets;`)
		if err != nil {
			t.Fatalf("Error Removing data table from data tsore :%v\n", err)
		}
		//close database connection
		if err = db.Close(); err != nil {
			t.Errorf("Error closing database file:%v\n", err)
		}
		// close tmp File
		if err = tmpFile.Close(); err != nil {
			t.Errorf("Error closing temporary file:%v\n", err)
		}
		//remove temp file
		if err = os.Remove(tmpFile.Name()); err != nil {
			t.Errorf("Error removing temporary file: %v\n", err)
		}
	}

	fmt.Println("Temporary file name:", tmpFile.Name())
	return db, cleanUp

}

func removeDb(filePath string, t *testing.T) {

	_, err := os.Stat(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			t.Error("MockDb does not exist.")
		} else {
			t.Errorf("Error checking mockDb existence: %s", err)
		}
		// No need to proceed with deletion if file doesn't exist
		return
	}

	err = os.Remove(filePath)
	if err != nil {
		t.Errorf("Error deleting mockDb: %s", err)
	} else {
		t.Error("mockDb deleted successfully.")
	}

}

func Db(lg logger.ILogger) *sql.DB {

	db, err := sql.Open("sqlite", "/snippetbox.sqlite")
	if err != nil {
		lg.Fatal("Error connecting to the database:", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		lg.Fatal(err.Error())
	}
	return db
}
