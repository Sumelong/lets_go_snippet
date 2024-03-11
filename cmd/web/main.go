package main

import (
	"snippetbox/pkg"
)

func main() {
	// set up configuration
	app := pkg.NewApp(pkg.EnvInstanceDev)
	app.Name("Snippet Box").
		Logger(pkg.LogInstanceStdLogger).
		WebServerAddress("", "").
		WebServer(pkg.ServerInstanceMux)
	app.Run()

}
