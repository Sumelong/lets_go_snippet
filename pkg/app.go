package pkg

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"snippetbox/storing"
	"time"
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

	loggerInstance int
	Logging        Logger

	StoreInstance int
	Storage       *sql.DB

	port      string
	host      string
	staticDir string

	envInstance    int
	serverInstance int
	webServer      IServer

	ctxTimeout time.Duration

	prodErrLogFile  string
	prodInfoLogFile string
}

func NewApp(envInstance int) *App {
	var (
		p, h, s string
	)

	flag.StringVar(&p, "port", "4000", "HTTP network address")
	flag.StringVar(&h, "host", "http://localhost", "HTTP network address")
	flag.StringVar(&s, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	return &App{
		port:        p,
		host:        h,
		staticDir:   s,
		envInstance: envInstance,
	}

}

func (a *App) Logger(instance int) *App {

	a.loggerInstance = instance

	//set names of production logging files
	// app will create a director in root call logs
	// and create your set files in the directory
	a.prodInfoLogFile = "logInfo.log"
	a.prodErrLogFile = "logErr.log"

	lg, errs := NewLoggerFactory(a)

	a.Logging = lg
	a.err = errs

	a.Logging.Info("app Logging configuration successful")

	return a
}

func (a *App) Name(name string) *App {

	switch name {
	case "":
		a.name = os.Getenv("APP_NAME")
	default:
		a.name = name
	}
	a.Logging.Info("app name configuration successful")
	return a
}

func (a *App) Store(instance int) *App {
	a.StoreInstance = instance
	a.Storage = storing.NewStoreFactory(a)
	return a
}

func (a *App) Migrate() {
	a.Logging.Info("beginning migration")
	storing.RunMigration(a)
	a.Logging.Info("migration completed")
}

func (a *App) WebServerAddress(p string, h string) *App {

	//set port and host of server

	//check if null and return appConfig to use default value
	if p == "" || h == "" {
		a.Logging.Info(fmt.Sprintf("set port-%s and host-%s", p, h))
		return a
	}
	//if not null use provided value
	a.port = p
	a.host = h
	a.Logging.Info(fmt.Sprintf("set port-%s and host-%s", p, h))
	a.Logging.Info("app address set")
	return a

}

func (a *App) WebServer(serverInstance int) *App {

	a.serverInstance = serverInstance
	// assign server from factory to app server
	a.webServer, a.err = NewServerFactory(a)

	a.Logging.Info("app server configured", nil)

	return a
}

func (a *App) Run() {

	if a.err != nil {
		log.Fatal(a.err)
	}

	a.Logging.Info("app begin successful")
	a.Logging.Info("starting app :%s", a.name)

	err := a.webServer.Begin()
	if err != nil {
		log.Fatal(err)
	}

}
