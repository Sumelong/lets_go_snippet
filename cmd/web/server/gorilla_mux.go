package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"snippetbox/cmd/web/handlers"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
	"time"

	"github.com/gorilla/mux"
)

type Gorilla struct {
	router *mux.Router
	handle *handlers.Handle
	logger logger.Logger
	addr   string
}

func NewGorillaMux(lg logger.Logger, addr string, snippet models.ISnippet) (*Gorilla, error) {

	h, err := handlers.NewHandle(snippet, lg)
	if err != nil {
		return nil, err
	}

	// return server
	return &Gorilla{
		router: mux.NewRouter(),
		logger: lg,
		addr:   addr,
		handle: h,
	}, nil
}

func (s *Gorilla) routes() http.Handler {
	//flag.StringVar(&m.app.port, "addr", "4000", "HTTP network address")
	//flag.Parse()

	dir := filepath.Join(".", "ui", "static")
	fileServer := http.FileServer(http.Dir(dir))
	s.router.Handle("/static", http.NotFoundHandler())
	s.router.Handle("/static/", http.StripPrefix("/static", fileServer))

	//s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	s.router.HandleFunc("/health", s.healthCheckerHandler) //.Methods(http.MethodGet)
	s.router.HandleFunc("/", s.homeHandler)
	s.router.HandleFunc("/snippet/create", s.createSnippetHandler)
	s.router.HandleFunc("/snippet/create", s.createSnippetFormHandler)
	s.router.HandleFunc("/snippet/{id:[0-9]+}", s.showSnippetHandler)

	// Wrap the existing chain with the logRequest middleware.
	return s.handle.RecoverPanic(s.handle.LogRequest(s.handle.SecureHeaders(s.router)))

}

func (s *Gorilla) healthCheckerHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.HealthChecker(w, r)
}

func (s *Gorilla) homeHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.Home(w, r)
}

func (s *Gorilla) showSnippetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	q := r.URL.Query()
	q.Add("snippet_id", vars["id"])
	r.URL.RawQuery = q.Encode()
	s.handle.ShowSnippet(w, r)
}

func (s *Gorilla) createSnippetHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.CreateSnippet(w, r)
}

func (s *Gorilla) createSnippetFormHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.CreateSnippetForm(w, r)
}

func (s *Gorilla) Begin() error {

	//set handle
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
