package store

import (
	"database/sql"
	"flag"
	_ "modernc.org/sqlite"
	"snippetbox/pkg/logger"
	"snippetbox/pkg/services"
)

func NewStoreSqlite(lg *logger.Logger) *sql.DB {

	var dsn string
	conn, err := services.FindFile("snippetbox.sqlite") // Load variables from .env file
	if err != nil {
		lg.Error("could not locate snippetbox store:", err)
		return nil
	}
	lg.Info("snippetbox store located successfully:")

	err = services.LoadEnv()
	if err != nil {
		lg.Error("Error loading .env file:", err)
		return nil
	}
	lg.Info(".env loaded successfully")

	flag.StringVar(&dsn, "pg-dsn", conn, "postgres data store")
	flag.Parse()

	//fmt.Println(dns)
	lg.Info("sqlite dsn configuration successful")

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		lg.Fatal("Error connecting to the database:", err)
		//log.Fatal("Error connecting to the database:", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		lg.Fatal(err.Error())
		//log.Fatalln(err)
	}

	return db
}
