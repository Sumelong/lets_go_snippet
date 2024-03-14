package test

import (
	"database/sql"
	"snippetbox/pkg/logger"
)

func db(lg logger.ILogger) *sql.DB {
	//c := NewConfigSqlite()

	//var dsn = fmt.Sprintf("%s", "pkg/domain/sqlite/snippetbox.sqlite")

	//fmt.Println(dns)
	//lg.Info("dns configuration successful")
	db, err := sql.Open("sqlite", "/snippetbox.sqlite")
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
