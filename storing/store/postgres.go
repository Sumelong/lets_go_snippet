package store

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"snippetbox/pkg/logger"
	"snippetbox/pkg/services"
)

func NewStorePostgres(lg *logger.Logger) *sql.DB {
	err := services.LoadEnv() // Load variables from .env file
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return nil
	}

	c := NewConfigPostgres()

	var dsn string
	conn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Host,
		c.Port,
		c.User,
		c.Database,
		c.Password,
	)
	flag.StringVar(&dsn, "pg-dsn", conn, "postgres data store")
	flag.Parse()

	if c.Host == "" {
		lg.Fatal("failed to configure store")
	}

	//fmt.Println(dns)
	lg.Info("postgres dsn configuration successful")
	db, err := sql.Open("postgres", dsn)
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
