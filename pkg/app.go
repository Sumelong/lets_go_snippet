package pkg

import (
	"errors"
	"flag"
	"log"
)

const (
	EnvInstanceDev int = iota
	EnvInstanceProd
)

var (
	ErrUnsupportedEnv = errors.New("unsupported environment")
)

type IApp interface {
	Run()
}

type App struct {
	name      string
	Port      string
	Host      string
	StaticDir string

	envInstance    int
	serverInstance int

	loggerInstance int
	logger         ILogger

	prodErrLogFile  string
	prodInfoLogFile string
}

func NewApp(envInstance int) *App {
	var (
		p, h, s string
	)

	flag.StringVar(&p, "port", ":4000", "HTTP network address")
	flag.StringVar(&h, "host", "http://localhost", "HTTP network address")
	flag.StringVar(&s, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	return &App{
		Port:        p,
		Host:        h,
		StaticDir:   s,
		envInstance: envInstance,
	}

}

func (a *App) Name(name string) *App {
	a.name = name
	return a
}

func (a *App) Logger(li int) *App {
	a.loggerInstance = li

	l, err := NewLoggerFactory(a)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	a.logger = l

	a.prodInfoLogFile = "logInfo.log"
	a.prodErrLogFile = "logErr.log"
	return a
}

func (a *App) WebServer(p string, h string) *App {
	a.Port = p
	a.Host = h
	return a
}

func (a *App) Run() {

	srv, err := NewServerMux(a)
	if err != nil {
		log.Fatalln(err)
	}

	srv.Run()

}
