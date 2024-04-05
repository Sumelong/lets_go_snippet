package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"snippetbox/pkg/domain/adapters/persister"
	"snippetbox/pkg/domain/ports"
	"snippetbox/pkg/services"
	"time"

	"github.com/golangcollege/sessions"
	"snippetbox/cmd/web/server"
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

type Application struct {
	name string
	err  error

	Logger            logger.ILogger
	storeConn         *sql.DB
	UserRepository    *ports.IUserRepository
	SnippetRepository *ports.ISnippetRepository
	Session           *sessions.Session
	secret            *string

	addr      string
	staticDir string

	webServer server.IServer

	ctxTimeout time.Duration

	prodErrLogFile  string
	prodInfoLogFile string
	envInstance     int
	storeInstance   int
}

func NewApplication(envInstance int) Application {
	err := services.LoadEnv() // Load variables from .env file
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return Application{}
	}

	addr := flag.String("port", "4000", "HTTP network address")
	dir := flag.String("static-dir", "./ui/static", "Path to static assets")
	secret := flag.String("secret", os.Getenv("SECRETE"), "Secret key")
	flag.Parse()

	return Application{
		addr:        *addr,
		staticDir:   *dir,
		envInstance: envInstance,
		secret:      secret,
	}

}

func (a Application) Name(name string) Application {

	switch name {
	case "":
		a.name = os.Getenv("APP_NAME")
	default:
		a.name = name
	}
	return a
}

func (a Application) Logging(logInstance int) Application {

	//set names of production logging files
	// app will create a director in root call logs
	// and create your set files in the directory
	infoLogFile := "logInfo.log"
	errLogFile := "logErr.log"

	a.Logger, a.err = logger.NewLoggerFactory(a.envInstance, logInstance, errLogFile, infoLogFile)
	/*
		a.logger = lg
		a.err = errs*/

	a.Logger.Info("app logger configuration successful\n")

	return a
}

func (a Application) Storing(storeInstance int) Application {

	a.storeInstance = storeInstance
	s := store.NewStoreFactory(storeInstance, &a.Logger)
	a.storeConn = s
	return a
}

func (a Application) Repository() Application {

	ur, sr, err := persister.NewRepositoryFactory(
		a.storeInstance,
		&a.Logger,
		a.storeConn,
	)
	a.UserRepository = ur
	a.SnippetRepository = sr

	a.err = err
	return a
}

func (a Application) Migrate() {

	a.Logger.Info("beginning migration")
	store.RunMigration(a.storeInstance, a.storeConn, &a.Logger)
	a.Logger.Info("migration completed")
}

func (a Application) WebServerAddress(addr *string) Application {

	//check if null and return appConfig to use default value
	if addr != nil {
		a.Logger.Info(fmt.Sprintf("set port-%s", *addr))
		//if not null use provided value
		a.addr = *addr
		a.Logger.Info(fmt.Sprintf("set port-%s", *addr))
		return a
	}
	return a

}

func (a Application) WebServer(serverInstance int) Application {

	// Use the sessions.New() function to initialize a new session manager,
	// passing in the secret key as the parameter. Then we configure it so
	// sessions always expires after 12 hours.
	session := sessions.New([]byte(*a.secret))
	session.Lifetime = 12 * time.Hour
	session.SameSite = http.SameSiteStrictMode
	a.Session = session

	staticFileDir, err := services.FindFile("ui/html")
	if err != nil {
		a.Logger.Error("static files not found %s", err)
	}

	// assign server from factory to app server
	srv, err := server.NewServerFactory(
		serverInstance,
		&a.Logger,
		a.addr,
		a.UserRepository,
		a.SnippetRepository,
		a.Session,
		staticFileDir,
	)

	a.webServer = srv
	a.err = err

	a.Logger.Info("webserver %d configured", serverInstance)

	return a
}

func (a Application) Run() {

	if a.err != nil {
		a.Logger.Fatal(a.err.Error(), a.err)
	}

	a.Logger.Info("app configuration successful")
	a.Logger.Info("starting app :%s", a.name)

	err := a.webServer.Begin()
	if err != nil {
		a.Logger.Fatal(err.Error(), err)
	}

}
