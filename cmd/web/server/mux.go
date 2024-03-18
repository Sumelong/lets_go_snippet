package server

import (
	"fmt"
	"net/http"
	"snippetbox/cmd/web/controller"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
	"time"
)

type Mux struct {
	router     *http.ServeMux
	controller *controller.Controller
	logger     logger.Logger
	addr       string
}

func NewMux(lg logger.Logger, addr string, snippet models.ISnippet) (*Mux, error) {

	handler, err := controller.NewController(snippet, lg)
	if err != nil {
		return nil, err
	}

	// return server
	return &Mux{
		router:     http.NewServeMux(),
		logger:     lg,
		addr:       addr,
		controller: handler,
	}, nil
}

func (s *Mux) routes() http.Handler {
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
	s.router.Handle("/static", http.NotFoundHandler())
	s.router.Handle("/static/", http.StripPrefix("/static", fileServer))

	s.router.HandleFunc("/health", s.controller.HealthChecker)
	s.router.HandleFunc("/", s.controller.Home)
	s.router.HandleFunc("/snippet/", s.controller.ShowSnippet)
	s.router.HandleFunc("/snippet/create", s.controller.CreateSnippet)

	// Wrap the existing chain with the logRequest middleware.
	return s.controller.RecoverPanic(s.controller.LogRequest(s.controller.SecureHeaders(s.router)))

}

func (s *Mux) Begin() error {

	//set controller
	//routes(s)

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:     fmt.Sprintf(":%s", s.addr),
		ErrorLog: s.logger.ErrLog,
		Handler:  s.routes(),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.logger.Info(fmt.Sprintf("server running on port:%s", s.addr))

	// Call the ListenAndServe() method on our new http.Server struct.
	err := srv.ListenAndServe()

	//log error from serve
	s.logger.Error(err.Error(), err)

	return err
}
