-- Create a new UTF-8 'snippetbox' database.
create table snippets(
                         id      INTEGER not null primary key autoincrement,
                         title   VARCHAR(100)   not null,
                         content TEXT          not null,
                         created TIMESTAMP default CURRENT_TIMESTAMP not null,
                         expires TIMESTAMP    not null);
create index idx_snippets_created on snippets (created);

-- add value to snippets
insert into snippets(title,content,created,expires)
values ('test title 1','test content 1',current_timestamp,date(current_timestamp+'2 days') );

insert into snippets(title,content,created,expires)
values ('test title 2','test content 2',current_timestamp,date(current_timestamp+'2 days') );

insert into snippets(title,content,created,expires)
values ('test title 3','test content 3',current_timestamp,date(current_timestamp+'2 days') );

-- create usr table
create table users(
                      id  INTEGER  not null primary key autoincrement,
                      name            VARCHAR(255)         not null,
                      email           VARCHAR(255)         not null
                          constraint users_uc_email unique,
                      hashed_password CHAR(60)             not null,
                      created         DATETIME             not null,
                      active          BOOLEAN default TRUE not null);
create index users_id_index on users (id);
-- add value to user table
INSERT INTO users (name, email, hashed_password, created)
VALUES ('Alice Jones','alice@example.com','$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG','2018-12-23 17:25:22');