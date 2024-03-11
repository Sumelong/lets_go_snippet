package migrate

import (
	"snippetbox/pkg"
	"snippetbox/storing"
)

func main() {
	// set up configuration
	app := pkg.NewApp(pkg.EnvInstanceDev)
	app.Name("Snippet Box").
		Logger(pkg.LogInstanceStdLogger).
		Store(storing.StoreInstancePostgres).
		Migrate()

}
