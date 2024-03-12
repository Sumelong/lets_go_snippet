package main

import (
	"snippetbox/pkg/logger"
	"snippetbox/storing/store"
)

func main() {
	// set up configuration
	app := NewApp(EnvInstanceDev)
	app.Name("Snippet-Migration").
		Logging(logger.LogInstanceStdLogger).
		Storing(store.StorageInstancePostgres).Migrate()

}
