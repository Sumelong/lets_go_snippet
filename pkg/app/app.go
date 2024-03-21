package app

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
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
)

const (
	EnvInstanceDev int = iota
	EnvInstanceProd
)

var (
	ErrUnsupportedEnv = errors.New("unsupported environment")
)

type App struct {
	Name string
	Err  error

	Logger  logger.Logger
	Store   *sql.DB
	Snippet *models.ISnippet
	Session *sessions.Session
	Secret  *string

	Addr      string
	StaticDir string

	WebServer server.IServer

	CtxTimeout time.Duration

	ProdErrLogFile  string
	ProdInfoLogFile string
	EnvInstance     int
	StoreInstance   int
}

func New(envInstance int) App {
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
		Addr:        *addr,
		StaticDir:   *dir,
		EnvInstance: envInstance,
		Secret:      secret,
	}
}

type IApp interface {
	Name(string) IApp
	Logger(int) IApp
	Store(int) IApp
	Model() IApp
	Migrate()
	WebServerAddress(*string) IApp
	WebServer(int) IApp
	Run()
}
