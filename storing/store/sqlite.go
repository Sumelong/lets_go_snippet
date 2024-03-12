package store

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"snippetbox/pkg"
)

func NewStoreSqlite(lg pkg.Logger) *sql.DB {
	c := NewConfigSqlite()

	var dns = fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Host,
		c.Port,
		c.User,
		c.Database,
		c.Password,
	)

	//fmt.Println(dns)
	lg.Info("dns configuration successful")
	db, err := sql.Open(c.Driver, dns)
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
