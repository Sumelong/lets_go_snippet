package main

import (
	"snippetbox/pkg"
	"snippetbox/storing/store"
)

func main() {
	// set up configuration
	app := NewApp(EnvInstanceDev)
	app.Name("Snippet-Migration").
		Logging(pkg.LogInstanceStdLogger).
		Storing(store.StorageInstancePostgres).Migrate()

}
