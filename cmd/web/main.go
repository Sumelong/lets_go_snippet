package main

import (
	"snippetbox/pkg/logger"
	"snippetbox/pkg/server"
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
		WebServer(server.ServerInstanceMux).Run()
	//app.Run()

}
