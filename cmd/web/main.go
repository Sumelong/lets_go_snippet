package main

import (
	"snippetbox/pkg"
)

func main() {
	// set up configuration
	app := NewApp(EnvInstanceDev)
	app.Name("Snippet Box").
		Logging(pkg.LogInstanceStdLogger).
		WebServerAddress("").
		WebServer(pkg.ServerInstanceMux).Run()
	//app.Run()

}
