package server

import (
	"fmt"
	"github.com/golangcollege/sessions"
	"github.com/gorilla/pat"
	"github.com/justinas/alice"
	"net/http"
	"path/filepath"
	"snippetbox/cmd/web/handlers"
	"snippetbox/pkg/domain/ports"
	"snippetbox/pkg/logger"
	"time"
)

type GorillaPat struct {
	router *pat.Router
	handle *handlers.Handle
	logger logger.ILogger
	addr   string
}

func NewGorillaPat(
	lg *logger.ILogger,
	addr string,
	user *ports.IUserRepository,
	snippet *ports.ISnippetRepository,
	session *sessions.Session,
	staticFileDir string,
) (*GorillaPat, error) {

	h, err := handlers.NewHandle(user, snippet, lg, session, staticFileDir)
	if err != nil {
		return nil, err
	}

	// return server
	return &GorillaPat{
		router: pat.New(),
		logger: *lg,
		addr:   addr,
		handle: h,
	}, nil
}

func (s *GorillaPat) routes() http.Handler {
	//flag.StringVar(&m.app.port, "addr", "4000", "HTTP network address")
	//flag.Parse()

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(s.handle.RecoverPanic, s.handle.LogRequest, s.handle.SecureHeaders)

	dir := filepath.Join(".", "ui", "static")
	fileServer := http.FileServer(http.Dir(dir))
	http.StripPrefix("/static", fileServer)
	s.router.Handle("/static/", http.StripPrefix("/static", fileServer))

	//s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	s.router.Get("/health", s.healthCheckerHandler)
	s.router.Get("/", s.homeHandler)
	s.router.Get("/snippet/{id}", s.showSnippetHandler)
	s.router.Get("/snippet/remove/{id}", s.removeSnippetHandler)
	s.router.Get("/snippet/create", s.createSnippetHandler)
	s.router.Post("/snippet/create", s.createSnippetFormHandler)

	return standardMiddleware.Then(s.router)
}

func (s *GorillaPat) healthCheckerHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.HealthChecker(w, r)
}

func (s *GorillaPat) homeHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.Home(w, r)
}

func (s *GorillaPat) showSnippetHandler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	q.Add("snippet_id", q.Get(":id"))
	r.URL.RawQuery = q.Encode()
	s.handle.ShowSnippet(w, r)
}

func (s *GorillaPat) removeSnippetHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("snippet_id", q.Get(":id"))
	r.URL.RawQuery = q.Encode()
	s.handle.RemoveSnippet(w, r)
}

func (s *GorillaPat) createSnippetHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.CreateSnippet(w, r)
}

func (s *GorillaPat) createSnippetFormHandler(w http.ResponseWriter, r *http.Request) {
	s.handle.CreateSnippetForm(w, r)
}

func (s *GorillaPat) Begin() error {

	//set handle
	//routes(s)

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", s.addr),
		//ErrorLog: s.logger.ErrLog,
		Handler: s.routes(),
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
