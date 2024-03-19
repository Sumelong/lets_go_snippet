package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"snippetbox/cmd/web/handlers"
	"time"

	"github.com/justinas/alice"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
)

type GoMux struct {
	router *http.ServeMux
	handle *handlers.Handle
	logger logger.Logger
	addr   string
}

func NewGoMux(lg logger.Logger, addr string, snippet models.ISnippet) (*GoMux, error) {

	c, err := handlers.NewHandle(snippet, lg)
	if err != nil {
		return nil, err
	}

	// return server
	return &GoMux{
		router: http.NewServeMux(),
		logger: lg,
		addr:   addr,
		handle: c,
	}, nil
}

func (s *GoMux) routes() http.Handler {
	//flag.StringVar(&m.app.port, "addr", "4000", "HTTP network address")
	//flag.Parse()

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(s.handle.RecoverPanic, s.handle.LogRequest, s.handle.SecureHeaders)

	dir := filepath.Join(".", "ui", "static")
	fileServer := http.FileServer(http.Dir(dir))
	s.router.Handle("/static", http.NotFoundHandler())
	s.router.Handle("/static/", http.StripPrefix("/static", fileServer))

	//s.router.PathPrefix("/static/").Handle(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	s.router.HandleFunc("/health", s.healthCheckerHandler)
	s.router.HandleFunc("/", s.homeHandler)
	s.router.HandleFunc("/snippet/{id}", s.showSnippetHandler)
	s.router.HandleFunc("/snippet/create", s.createSnippetHandler)
	s.router.HandleFunc("/snippet/create", s.createSnippetFormHandler)

	// Return the 'standard' middleware chain followed by the servemux router.
	return standardMiddleware.Then(s.router)

}

func (s *GoMux) healthCheckerHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.HealthChecker(w, r)
}

func (s *GoMux) homeHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.Home(w, r)
}

func (s *GoMux) showSnippetHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("snippet_id", q.Get("id"))
	r.URL.RawQuery = q.Encode()
	s.handle.ShowSnippet(w, r)
}

func (s *GoMux) createSnippetHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.CreateSnippet(w, r)
}

func (s *GoMux) createSnippetFormHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.CreateSnippetForm(w, r)
}

func (s *GoMux) Begin() error {

	//set handle
	//routes(s)

	// Initialize a new http.Server struct. We set the Addr and Handle fields so
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
