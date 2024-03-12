package server

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
)

type ServerMux struct {
	mx   *http.ServeMux
	hdl  *Handlers
	lg   logger.Logger
	addr string
}

func NewServerMux(lg logger.Logger, addr string, snippet models.ISnippet) (*ServerMux, error) {

	// return server
	return &ServerMux{
		mx:   http.NewServeMux(),
		lg:   lg,
		addr: addr,
		hdl:  NewHandler(snippet, lg),
	}, nil
}

func (m ServerMux) setHandler(r *http.ServeMux) {

	r.HandleFunc("/health", m.hdl.HealthChecker)
	r.HandleFunc("/", m.hdl.Home)
	r.HandleFunc("/snippet", m.hdl.ShowSnippet)
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
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:     fmt.Sprintf(":%s", m.addr),
		ErrorLog: m.lg.ErrLog,
		Handler:  m.mx,
	}

	m.lg.Info(fmt.Sprintf("server running on port:%s", m.addr))

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

var ErrUnsupportedServer = errors.New("unsupported server")

func NewServerFactory(serverInstance int, lg logger.Logger, addr string, snippet models.ISnippet) (IServer, error) {

	switch serverInstance {
	case ServerInstanceCustom:
		return nil, ErrUnsupportedServer
	case ServerInstanceMux:
		return NewServerMux(lg, addr, snippet)
	default:
		return nil, ErrUnsupportedServer
	}
}
