package migration

import (
	"database/sql"
	"log"
	"snippetbox/pkg"
)

func NewPostgresMigration(app pkg.App) {

	db := app.Storage
	lg := app.Logging
	defer db.Close()

	/*// Create the database
	// Replace "create_db.sh" with the actual path to your script
	cmd := exec.Command("sh", "create_db.sh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		lg.Error("Error creating database:", err)
		lg.Info(string(output))
		return
	}
	fmt.Println("Database created successfully!")*/

	// Create the snippet table
	_, err := db.Exec(`CREATE TABLE snippets (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT,
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires TIMESTAMP);
		CREATE INDEX idx_snippets_created ON snippets(created);`)
	if err != nil {
		lg.Fatal("Error creating table:", err)

	}
}

func addSnippets(db *sql.DB, lg pkg.Logger) {
	// Create the snippet table
	_, err := db.Exec(`
INSERT INTO snippets (title, content, created, expires)
VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP + INTERVAL '1 year');

-- Add some dummy records (which we'll use in the next couple of chapters).
INSERT INTO snippets (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP + INTERVAL '365 days'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP + INTERVAL '365 days'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP + INTERVAL '7 days'
);




`)

	if err != nil {
		lg.Fatal(err.Error())
	}
}

func createUser(db *sql.DB, lg pkg.Logger) {
	_, err := db.Exec(`
	-- Create user 'web' with password 'pass'
	CREATE USER web WITH  ENCRYPTED PASSWORD 'snippets@pass';

	-- Grant privileges on the 'snippetbox' schema to user 'web'
	GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA snippetbox TO web;

	-- Alter the password for user 'web'
	ALTER USER web WITH ENCRYPTED PASSWORD 'new_password';
	`)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
