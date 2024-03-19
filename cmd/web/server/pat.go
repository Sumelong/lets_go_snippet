package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"snippetbox/cmd/web/handlers"
	"time"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
)

type Pat struct {
	router *pat.PatternServeMux
	handle *handlers.Handle
	logger logger.Logger
	addr   string
}

func NewPat(lg logger.Logger, addr string, snippet models.ISnippet) (*Pat, error) {

	h, err := handlers.NewHandle(snippet, lg)
	if err != nil {
		return nil, err
	}

	// return server
	return &Pat{
		router: pat.New(),
		logger: lg,
		addr:   addr,
		handle: h,
	}, nil
}

func (s *Pat) routes() http.Handler {
	//flag.StringVar(&m.app.port, "addr", "4000", "HTTP network address")
	//flag.Parse()

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(s.handle.RecoverPanic, s.handle.LogRequest, s.handle.SecureHeaders)

	// Create a file server which serves files out of the "./ui/static" directory.
	//	// Note that the path given to the http.Dir function is relative to the project
	//	// directory root.

	dir := filepath.Join(".", "ui", "static")
	fileServer := http.FileServer(http.Dir(dir))
	s.router.Get("/static", http.NotFoundHandler())
	s.router.Get("/static/", http.StripPrefix("/static", fileServer))

	//s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	s.router.Get("/health", http.HandlerFunc(s.healthCheckerHandler))
	s.router.Get("/", http.HandlerFunc(s.homeHandler))
	s.router.Get("/snippet/create", http.HandlerFunc(s.createSnippetFormHandler))
	s.router.Post("/snippet/create", http.HandlerFunc(s.createSnippetHandler))
	s.router.Get("/snippet/:id", http.HandlerFunc(s.showSnippetHandler))

	// Return the 'standard' middleware chain followed by the servemux router.
	return standardMiddleware.Then(s.router)
}

func (s *Pat) healthCheckerHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.HealthChecker(w, r)
}

func (s *Pat) homeHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.Home(w, r)
}

func (s *Pat) showSnippetHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("snippet_id", q.Get(":id"))
	r.URL.RawQuery = q.Encode()
	s.handle.ShowSnippet(w, r)
}

func (s *Pat) createSnippetHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.CreateSnippet(w, r)
}

func (s *Pat) createSnippetFormHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.CreateSnippetForm(w, r)
}

func (s *Pat) Begin() error {

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
