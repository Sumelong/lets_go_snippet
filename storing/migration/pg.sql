-- Create a new UTF-8 'snippetbox' database.
CREATE DATABASE snippetbox WITH ENCODING='UTF8';

-- Connect to the 'snippetbox' database.
\c snippetbox


-- Create user 'web' with password 'pass'
CREATE USER web WITH  ENCRYPTED PASSWORD 'web@pass';

-- Grant privileges on the 'snippetbox' schema to user 'web'
GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA snippetbox TO web;

-- Alter the password for user 'web'
ALTER USER web WITH PASSWORD 'new_password';



-- Create a 'snippets' table.
CREATE TABLE snippets (
                          id SERIAL PRIMARY KEY,
                          title VARCHAR(100) NOT NULL,
                          content TEXT NOT NULL,
                          created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          expires TIMESTAMP NOT NULL
);
-- Add an index on the created column.
CREATE INDEX idx_snippets_created ON snippets(created);



-- Add some dummy records (which we'll use in the next couple of chapters).
INSERT INTO snippets (title, content, created, expires)
VALUES (
       'An old silent pond',
       'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
       CURRENT_TIMESTAMP,
       CURRENT_TIMESTAMP + INTERVAL '365 days'
   );

INSERT INTO snippets (title, content, created, expires)
VALUES (
       'Over the wintry forest',
       'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
       CURRENT_TIMESTAMP,
       CURRENT_TIMESTAMP + INTERVAL '365 days'
   );

INSERT INTO snippets (title, content, created, expires)
VALUES (
       'First autumn morning',
       'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
       CURRENT_TIMESTAMP,
       CURRENT_TIMESTAMP + INTERVAL '7 days'
   );
