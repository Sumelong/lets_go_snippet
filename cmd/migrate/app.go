package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"snippetbox/pkg"
	"snippetbox/storing/store"
)

const (
	EnvInstanceDev int = iota
	EnvInstanceProd
)

var (
	ErrUnsupportedEnv = errors.New("unsupported environment")
)

type App struct {
	name string
	err  error

	logger pkg.Logger
	store  *sql.DB

	prodErrLogFile  string
	prodInfoLogFile string
	envInstance     int
	storeInstance   int
}

func NewApp(envInstance int) *App {
	err := godotenv.Load() // Load variables from .env file
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return &App{}
	}

	/*addr := flag.String("port", "4000", "HTTP network address")
	dir := flag.String("static-dir", "./ui/static", "Path to static assets")
	flag.Parse()*/

	return &App{
		envInstance: envInstance,
	}

}

func (a App) Name(name string) App {

	switch name {
	case "":
		a.name = os.Getenv("APP_NAME")
	default:
		a.name = name
	}
	return a
}

func (a App) Logging(logInstance int) App {

	//set names of production logging files
	// app will create a director in root call logs
	// and create your set files in the directory
	infoLogFile := "logInfo.log"
	errLogFile := "logErr.log"

	lg, errs := pkg.NewLoggerFactory(a.envInstance, logInstance, errLogFile, infoLogFile)

	a.logger = lg
	a.err = errs

	a.logger.Info("app logger configuration successful")

	return a
}

func (a App) Storing(storeInstance int) App {

	a.storeInstance = storeInstance
	a.store = store.NewStoreFactory(storeInstance, a.logger)
	return a
}

func (a App) Migrate() {

	a.logger.Info("beginning migration")
	store.RunMigration(a.storeInstance, a.store, a.logger)
	a.logger.Info("migration completed")
}
