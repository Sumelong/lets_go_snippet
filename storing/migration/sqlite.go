package migration

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
	"snippetbox/pkg/logger"
	"snippetbox/pkg/services"
	"testing"
)

func NewSqliteMigration(db *sql.DB, log *logger.ILogger) {

	//deference logger for usage
	lg := *log
	lg.Info("running migration on sqlite  data store ")

	// Create your table(s) here. This is just an example.

	//snippet table
	_, err := db.Exec(`
	create table snippets(
		id      INTEGER not null primary key autoincrement,
		title   VARCHAR(100)   not null,
		content TEXT          not null,
		created TIMESTAMP default CURRENT_TIMESTAMP not null,
	    expires TIMESTAMP    not null);
	create index idx_snippets_created on snippets (created);`)
	if err != nil {
		lg.Error("error creating snippets table:", err)
		return
	}
	//values
	_, err = db.Exec(`
	insert into snippets(title,content,created,expires)
	values ("test title 1","test content 1",current_timestamp,date(current_timestamp+'2 days') );

	insert into snippets(title,content,created,expires)
		values ("test title 2","test content 2",current_timestamp,date(current_timestamp+'2 days') );
	
	insert into snippets(title,content,created,expires)
			values ("test title 3","test content 3",current_timestamp,date(current_timestamp+'2 days') );
`)
	if err != nil {
		lg.Error("error adding values to snippets table:", err)
		return
	}

	// tenants table
	_, err = db.Exec(`
	create table tenants(
			id integer  not null  primary key autoincrement
			    constraint tenants_pk 
			    constraint tenants_pk_2 unique,
			created         DATETIME             not null,
			name   VARCHAR(255),
			code   varchar generated always as ('ES00' + id) 
			    stored not null,
			active BOOLEAN default true not null
	);`)
	if err != nil {
		lg.Error("error creating tenants table:", err)
		return
	}
	//values
	_, err = db.Exec(`
		insert into tenants(name,created,active) values ("tenant_1",current_timestamp,1);
		insert into tenants(name,created,active) values ("tenant_2",current_timestamp,1);
		insert into tenants(name,created,active) values ("tenant_3",current_timestamp,0);
		insert into tenants(name,created,active) values ("tenant_4",current_timestamp,1);
	`)
	if err != nil {
		lg.Error("error adding values to snippets table:", err)
		return
	}

	// tenants table
	_, err = db.Exec(`
	create table users(
			id  INTEGER  not null primary key autoincrement,
			name            VARCHAR(255)         not null,
			email           VARCHAR(255)         not null 
			    constraint users_uc_email unique,
			hashed_password CHAR(60)             not null,
			created         DATETIME             not null,
			active          BOOLEAN default TRUE not null);
	create index users_id_index on users (id);
   `)
	if err != nil {
		lg.Error("error creating users table:", err)
		return
	}
	//values
	_, err = db.Exec(`
		insert into users(name,email,hashed_password,created,active) 
		values ("admin","admin@email.com","xxxzzz",current_timestamp,1);

		insert into users(name,email,hashed_password,created,active) 
		values ("user_2","email2@example.com","xxxzzz",current_timestamp,1);

		insert into users(name,email,hashed_password,created,active) 
		values ("user_3","email3@example.com","xxxzzz",current_timestamp,1);
	`)
	if err != nil {
		lg.Error("error adding values to users table:", err)
		return
	}

}

func NewTestSqliteDB(t *testing.T) (*sql.DB, func()) {

	// Establish a sql.DB connection pool for our test database. Because our
	// setup and teardown scripts contains multiple SQL statements, we need
	// to use the `multiStatements=true` parameter in our DSN. This instructs
	// our MySQL database driver to support executing multiple SQL statements
	// in one db.Exec()` call.
	// Connect to the SQLite database
	db, err := sql.Open("sqlite", "test_snippetbox.db")
	if err != nil {
		t.Fatalf("Error connecting to the database: %v\n", err)
	}

	cleanUp := func() {
		//clean up function
		dbFile, err := services.FindFile("test_snippetbox.db")
		if err != nil {
			return
			//t.Fatalf("error locating test database: %v\n", err)
		}

		//remove database file
		if err = os.Remove(dbFile); err != nil {
			t.Fatalf("Error removing databse file: %v\n", err)
		}
	}
	cleanUp()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Error opening database: %v\n", err)
	}

	// Read the setup SQL script from file and execute the statements.
	setUpFile, err := services.FindFile("storing/migration/sqlite/setup.sql")
	if err != nil {
		t.Fatal("could not locate setup.sql script")
	}
	//script, err := os.ReadFile("./testdata/teardown.sql")
	setUpScript, err := os.ReadFile(setUpFile)
	//script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(setUpScript))

	if err != nil {
		t.Fatal(err)
	}

	// Return the connection pool and an anonymous function which reads and
	// executes the teardown script, and closes the connection pool.
	return db, func() {
		tearDownFile, err := services.FindFile("storing/migration/sqlite/teardown.sql")
		if err != nil {
			t.Fatal("could not locate teardown.sql script")
		}
		//script, err := os.ReadFile("./testdata/teardown.sql")
		tearDownScript, err := os.ReadFile(tearDownFile)
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(tearDownScript))

		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
}

func newSqliteMigration(log *logger.ILogger) (*sql.DB, func()) {
	lg := *log

	// Establish a sql.DB connection pool for our test database. Because our
	// setup and teardown scripts contains multiple SQL statements, we need
	// to use the `multiStatements=true` parameter in our DSN. This instructs
	// our Sqlite database driver to support executing multiple SQL statements
	// in one db.Exec()` call.
	// Connect to the SQLite database
	db, err := sql.Open("sqlite", "test_snippetbox.db")
	if err != nil {
		lg.Fatal("Error connecting to the database: %v\n", err)

	}
	err = db.Ping()
	if err != nil {
		lg.Fatal("Error opening database: %v\n", err)
	}

	// Read the setup SQL script from file and execute the statements.
	setUpFile, err := services.FindFile("storing/migration/sqlite/setup.sql")
	if err != nil {
		lg.Fatal("Error locating setup.sql script : %v", err)
	}
	//script, err := os.ReadFile("./testdata/teardown.sql")
	script, err := os.ReadFile(setUpFile)
	//script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		lg.Fatal("Error reading setup.sql script : %v", err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		lg.Fatal("Error executing setup.sql script : %v", err)
	}
	// Return the connection pool and an anonymous function which reads and
	// executes the teardown script, and closes the connection pool. We can
	// assign this anonymous function and call it later once our test has
	// completed.
	return db, func() {
		tearDownFile, err := services.FindFile("storing/migration/sqlite/setup.sql")
		if err != nil {
			lg.Fatal("Error locating teardown.sql script : %v", err)
		}
		//script, err := os.ReadFile("./testdata/teardown.sql")
		tearDownScript, err := os.ReadFile(tearDownFile)
		if err != nil {
			lg.Fatal("Error reading teardown.sql script : %v", err)
		}
		_, err = db.Exec(string(tearDownScript))
		if err != nil {
			lg.Fatal("Error executing teardown.sql script : %v", err)
		}
		db.Close()
	}
}
