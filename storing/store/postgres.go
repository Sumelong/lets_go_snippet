package store

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"snippetbox/pkg"
)

func NewStorePostgres(lg pkg.Logger) *sql.DB {
	c := NewConfigPostgres()

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
		lg.Fatal("Error connecting to the database:", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		lg.Debug(err.Error())
		//lg.Fatal(err.Error())
	}

	return db
}
