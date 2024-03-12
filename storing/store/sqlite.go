package store

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
	"snippetbox/pkg/logger"
)

func NewStoreSqlite(lg logger.Logger) *sql.DB {

	//c := NewConfigSqlite()

	//var dsn = fmt.Sprintf("%s", c.Database)

	//fmt.Println(dns)
	lg.Info("dns configuration successful")
	db, err := sql.Open("sqlite", "snippetbox.sqlite")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
