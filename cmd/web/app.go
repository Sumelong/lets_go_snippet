package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"snippetbox/pkg/services"
	"time"

	"github.com/golangcollege/sessions"
	"snippetbox/cmd/web/server"
	"snippetbox/pkg/domain"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
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

	logger  *logger.Logger
	store   *sql.DB
	snippet *models.ISnippet
	session *sessions.Session
	secret  *string

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
	err := services.LoadEnv() // Load variables from .env file
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return App{}
	}

	addr := flag.String("port", "4000", "HTTP network address")
	dir := flag.String("static-dir", "./ui/static", "Path to static assets")
	secret := flag.String("secret", os.Getenv("SECRETE"), "Secret key")
	flag.Parse()

	return App{
		addr:        *addr,
		staticDir:   *dir,
		envInstance: envInstance,
		secret:      secret,
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

	a.logger = &lg
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

	model, err := domain.NewSnippetsFactory(
		a.storeInstance,
		a.logger,
		a.store,
	)
	a.snippet = &model
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

	// Use the sessions.New() function to initialize a new session manager,
	// passing in the secret key as the parameter. Then we configure it so
	// sessions always expires after 12 hours.
	session := sessions.New([]byte(*a.secret))
	session.Lifetime = 12 * time.Hour
	a.session = session

	// assign server from factory to app server
	srv, err := server.NewServerFactory(
		serverInstance,
		a.logger,
		a.addr,
		a.snippet,
		session,
	)

	a.webServer = srv
	a.err = err

	a.logger.Info("webserver %d configured", serverInstance)

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
