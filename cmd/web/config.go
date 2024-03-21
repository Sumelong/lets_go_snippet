package main

import (
	"snippetbox/pkg/app"
)

type Config struct {
	app app.App
}

/*
func NewConfig(envInstance int) app.IApp {
	return Config{
		app: app.App{
			EnvInstance: envInstance,
		},
	}
}

func (c Config) Name(name string) app.IApp {

	switch name {
	case "":
		c.app.Name = os.Getenv("APP_NAME")
	default:
		c.app.Name = name
	}
	return c
}

func (c Config) Logger(logInstance int) app.IApp {

	//set names of production logging files
	// app will create a director in root call logs
	// and create your set files in the directory
	infoLogFile := "logInfo.log"
	errLogFile := "logErr.log"

	lg, errs := logger.NewLoggerFactory(
		c.app.EnvInstance,
		logInstance,
		errLogFile,
		infoLogFile,
	)

	c.app.Logger = lg
	c.app.Err = errs

	c.app.Logger.Info("app logger configuration successful\n")

	return c
}

func (c Config) Store(storeInstance int) app.IApp {

	c.app.StoreInstance = storeInstance
	s := store.NewStoreFactory(storeInstance, c.app.Logger)
	c.app.Store = s
	return c
}

func (c Config) Model() app.IApp {

	model, err := domain.NewSnippetsFactory(
		c.app.StoreInstance,
		c.app.Logger,
		c.app.Store,
	)
	c.app.Snippet = &model
	c.app.Err = err
	return c
}

func (c Config) Migrate() {

	c.app.Logger.Info("beginning migration")
	store.RunMigration(
		c.app.StoreInstance,
		c.app.Store,
		c.app.Logger,
	)
	c.app.Logger.Info("migration completed")
}

func (c Config) WebServerAddress(addr *string) app.IApp {

	//check if null and return appConfig to use default value
	if addr != nil {
		c.app.Logger.Info(fmt.Sprintf("set port-%s", *addr))
		//if not null use provided value
		c.app.Addr = *addr
		c.app.Logger.Info(fmt.Sprintf("set port-%s", *addr))
		return c
	}

	return c

}

func (c Config) WebServer(serverInstance int) app.IApp {

	// Use the sessions.New() function to initialize a new session manager,
	// passing in the secret key as the parameter. Then we configure it so
	// sessions always expires after 12 hours.
	session := sessions.New([]byte(*c.app.Secret))
	session.Lifetime = 12 * time.Hour
	c.app.Session = session

	// assign server from factory to app server
	srv, err := server.NewServerFactory(serverInstance, &c.app)

	c.app.WebServer = srv
	c.app.Err = err

	c.app.Logger.Info("webserver %d configured", serverInstance)

	return c
}

func (c Config) Run() {

	if c.app.Err != nil {
		c.app.Logger.Fatal(c.app.Err.Error(), c.app.Err)
	}

	c.app.Logger.Info("app configuration successful")
	c.app.Logger.Info("starting app :%s", c.app.Name)

	err := c.app.WebServer.Begin()
	if err != nil {
		c.app.Logger.Fatal(err.Error(), err)
	}

}
*/
