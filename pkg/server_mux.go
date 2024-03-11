package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
)

type ServerMux struct {
	mx  *http.ServeMux
	hdl Handlers
	lg  Logger
	app *App
}

func NewServerMux(a *App) (*ServerMux, error) {

	// return server
	return &ServerMux{
		mx:  http.NewServeMux(),
		lg:  a.Logging,
		app: a,
	}, nil
}

func (m ServerMux) setHandler(r *http.ServeMux) {

	r.HandleFunc("/health", m.hdl.HealthChecker)
	r.HandleFunc("/", m.hdl.Home)
	r.HandleFunc("/snippet/", m.hdl.ShowSnippet)
	r.HandleFunc("/snippet/create", m.hdl.CreateSnippet)
}

func (m ServerMux) Begin() error {

	//flag.StringVar(&m.app.port, "addr", "4000", "HTTP network address")
	//flag.Parse()

	// Create a file server which serves files out of the "./ui/static" directory.
	//	// Note that the path given to the http.Dir function is relative to the project
	//	// directory root.
	//******>>>fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	//******>>>mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	m.mx.Handle("/static", http.NotFoundHandler())
	m.mx.Handle("/static/", http.StripPrefix("/static", fileServer))

	//set handlers
	m.setHandler(m.mx)

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog Logging in
	// the event of any problems.
	srv := &http.Server{
		Addr:     fmt.Sprintf(":%s", m.app.port),
		ErrorLog: m.lg.ErrLog,
		Handler:  m.mx,
	}

	m.lg.Info("Starting server on %s", m.app.port)

	// Call the ListenAndServe() method on our new http.Server struct.
	err := srv.ListenAndServe()

	//log error from serve
	m.lg.Error(err.Error(), err)

	return err
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

type IServer interface {
	Begin() error
}

type port string

var ErrUnsupportedServer = errors.New("unsupported server")

func NewServerFactory(app *App) (IServer, error) {
	switch app.serverInstance {
	case ServerInstanceCustom:
		return nil, ErrUnsupportedServer
	case ServerInstanceMux:
		return NewServerMux(app)
	default:
		return nil, ErrUnsupportedServer
	}
}
