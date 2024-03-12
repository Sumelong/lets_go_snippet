package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"snippetbox/pkg/domain"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
	"snippetbox/pkg/server"
	"time"

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

	logger  logger.Logger
	store   *sql.DB
	snippet models.ISnippet

	addr      string
	staticDir string

	webServer server.IServer

	ctxTimeout time.Duration

	prodErrLogFile  string
	prodInfoLogFile string
	envInstance     int
	storeInstance   int
}

func NewApp(envInstance int) App {
	err := godotenv.Load() // Load variables from .env file
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return App{}
	}

	addr := flag.String("port", "4000", "HTTP network address")
	dir := flag.String("static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	return App{
		addr:        *addr,
		staticDir:   *dir,
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

	lg, errs := logger.NewLoggerFactory(a.envInstance, logInstance, errLogFile, infoLogFile)

	a.logger = lg
	a.err = errs

	a.logger.Info("app logger configuration successful\n")

	return a
}

func (a App) Storing(storeInstance int) App {

	a.storeInstance = storeInstance
	s := store.NewStoreFactory(storeInstance, a.logger)
	a.store = s
	return a
}

func (a App) Model() App {

	model, err := domain.NewSnippetsFactory(a.storeInstance, a.logger, a.store)
	a.snippet = model
	a.err = err
	return a
}

func (a App) Migrate() {

	a.logger.Info("beginning migration")
	store.RunMigration(a.storeInstance, a.store, a.logger)
	a.logger.Info("migration completed")
}

func (a App) WebServerAddress(addr *string) App {

	//check if null and return appConfig to use default value
	if addr != nil {
		a.logger.Info(fmt.Sprintf("set port-%s", *addr))
		//if not null use provided value
		a.addr = *addr
		a.logger.Info(fmt.Sprintf("set port-%s", *addr))
		return a
	}

	return a

}

func (a App) WebServer(serverInstance int) App {

	// assign server from factory to app server
	srv, err := server.NewServerFactory(serverInstance, a.logger, a.addr, a.snippet)

	a.webServer = srv
	a.err = err

	a.logger.Info("app server configured", nil)

	return a
}

func (a App) Run() {

	if a.err != nil {
		a.logger.Fatal(a.err.Error(), a.err)
	}

	a.logger.Info("app configuration successful")
	a.logger.Info("starting app :%s", a.name)

	err := a.webServer.Begin()
	if err != nil {
		a.logger.Fatal(err.Error(), err)
	}

}
