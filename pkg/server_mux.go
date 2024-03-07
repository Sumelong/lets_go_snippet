package pkg

import (
	"errors"
	"net/http"
	"path/filepath"
)

type ServerMux struct {
	mx     *http.ServeMux
	hdl    Handlers
	lg     *Logger
	config *App
}

func NewServerMux(a *App) (ServerMux, error) {
	//initial new logger
	lg, err := NewLoggerFactory(a)
	if err != nil {
		return ServerMux{}, err
	}

	// return server
	return ServerMux{
		mx:     http.NewServeMux(),
		lg:     a,
		config: a,
	}, nil
}

func (m ServerMux) setHandler(r *http.ServeMux) {

	r.HandleFunc("/health", HealthChecker)
	r.HandleFunc("/", Home)
	r.HandleFunc("/snippet/", ShowSnippet)
	r.HandleFunc("/snippet/create", CreateSnippet)
}

func (m ServerMux) Run() {

	// Create a file server which serves files out of the "./ui/static" directory.
	//	// Note that the path given to the http.Dir function is relative to the project
	//	// directory root.
	//******>>>fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	//******>>>mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	m.mx.Handle("/static", http.NotFoundHandler())
	m.mx.Handle("/static/", http.StripPrefix("/static", fileServer))

	m.setHandler(m.mx)
	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:     m.config.Port,
		ErrorLog: m.lg.ErrLog,
		Handler:  m.mx,
	}

	m.lg.InfoLog.Printf("Starting server on %s", m.config.Port)

	// Call the ListenAndServe() method on our new http.Server struct.
	err := srv.ListenAndServe()

	//log and panic of server error
	m.lg.ErrLog.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err = nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

//**************** SERVER FACTORY ****************///

const (
	ServerInstanceCustom int = iota
	ServerInstanceMux
)

var ErrUnsupportedServer = errors.New("unsupported server")

type appServer struct {
}

func (s appServer) NewServerFactory(app *App) (IApp, error) {
	switch app.serverInstance {
	case ServerInstanceCustom:
		return nil, ErrUnsupportedServer
	default:
		return NewServerMux(app)
	}
}
