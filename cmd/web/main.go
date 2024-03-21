package main

import (
	"snippetbox/cmd/web/server"
	"snippetbox/pkg/logger"
	"snippetbox/storing/store"
)

func main() {
	// set up configuration
	app := NewApp(EnvInstanceDev)
	app.Name("Snippet Box").
		Logging(logger.LogInstanceStdLogger).
		Storing(store.StorageInstanceSqlite).
		Model().
		WebServerAddress(nil).
		WebServer(server.ServeInstancePat).Run()
	//app.Run()

}

/*
func main() {
	// set up configuration
	c := NewConfig(EnvInstanceDev)
	c.Name("Snippet Box").
		Logger(logger.LogInstanceStdLogger).
		Store(store.StorageInstanceSqlite).
		Model().
		WebServerAddress(nil).
		WebServer(server.ServeInstancePat).Run()
	//app.Run()

}*/
