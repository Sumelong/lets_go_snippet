package main

import (
	"snippetbox/pkg"
)

func main() {
	// set up configuration
	app := pkg.NewApp(pkg.EnvInstanceProd)
	app.Name("Snippet Box").
		Logger(pkg.LogInstanceStdLogger).
		WebServer("", "")
	app.Run()

}
