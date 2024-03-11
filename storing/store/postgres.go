package store

import (
	"database/sql"
	"fmt"
	"log"
	"snippetbox/pkg"
	"snippetbox/storing"

	_ "github.com/lib/pq"
)

func NewStorePostgres(lg pkg.Logger) *sql.DB {
	c := storing.NewConfigPostgres()

	var dns = fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Host,
		c.Port,
		c.User,
		c.Database,
		c.Password,
	)

	//fmt.Println(dns)
	lg.Info("dsn configuration successful")
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
