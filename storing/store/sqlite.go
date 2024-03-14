package store

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"snippetbox/pkg/logger"
	"snippetbox/pkg/services"
)

func NewStoreSqlite(lg logger.Logger) *sql.DB {

	dns, err := services.FindFile("snippetbox.sqlite") // Load variables from .env file
	if err != nil {
		fmt.Println("could not locate snippetbox store:", err)
		return nil
	}
	err = services.LoadEnv()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return nil
	}

	//fmt.Println(dns)
	lg.Info("dsn configuration successful")

	db, err := sql.Open("sqlite", dns)
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
