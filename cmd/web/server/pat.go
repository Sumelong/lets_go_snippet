package server

import (
	"fmt"
	"github.com/golangcollege/sessions"
	"net/http"
	"path/filepath"
	"snippetbox/cmd/web/handlers"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
	"time"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

type Pat struct {
	router  *pat.PatternServeMux
	handle  *handlers.Handle
	logger  *logger.Logger
	addr    string
	session *sessions.Session
}

/*type Pat struct {
	router *pat.PatternServeMux
	handle *handlers.Handle
	app    *app.App
}*/

func NewPat(lg *logger.Logger, addr string, snippet *models.ISnippet, session *sessions.Session) (*Pat, error) {
	//func NewPat(app *app.App) (*Pat, error) {

	//h, err := handlers.NewHandle(snippet, lg)
	h, err := handlers.NewHandle(snippet, lg, session)
	if err != nil {
		return nil, err
	}

	// return server
	return &Pat{
		router:  pat.New(),
		logger:  lg,
		addr:    addr,
		handle:  h,
		session: session,
	}, nil

	/*// return server
	return &Pat{
		router: pat.New(),
		handle: h,
		app:    app,
	}, nil
	*/
}

func (s *Pat) routes() http.Handler {
	//flag.StringVar(&m.app.port, "addr", "4000", "HTTP network address")
	//flag.Parse()

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(s.handle.RecoverPanic, s.handle.LogRequest, s.handle.SecureHeaders)

	// Create a new middleware chain containing the middleware specific to
	// our dynamic application routes. For now, this chain will only contain
	// the session middleware but we'll add more to it later.
	dynamicMiddleware := alice.New(s.session.Enable)

	// Create a file server which serves files out of the "./ui/static" directory.
	//	// Note that the path given to the http.Dir function is relative to the project
	//	// directory root.

	dir := filepath.Join(".", "ui", "static")
	fileServer := http.FileServer(http.Dir(dir))
	s.router.Get("/static", http.NotFoundHandler())
	s.router.Get("/static/", http.StripPrefix("/static", fileServer))

	//s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	s.router.Get("/health", http.HandlerFunc(s.healthCheckerHandler))
	// Update these routes to use the new dynamic middleware chain followed
	// by the appropriate handler function.
	s.router.Get("/", dynamicMiddleware.ThenFunc(s.homeHandler))
	//s.router.Get("/", http.HandlerFunc(s.homeHandler))
	s.router.Get("/snippet/create", dynamicMiddleware.ThenFunc(s.createSnippetFormHandler))
	s.router.Post("/snippet/create", dynamicMiddleware.ThenFunc(s.createSnippetHandler))
	s.router.Post("/snippet/remove/:id", dynamicMiddleware.ThenFunc(s.removeSnippetHandler))
	s.router.Get("/snippet/:id", dynamicMiddleware.ThenFunc(s.showSnippetHandler))

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

func (s *Pat) removeSnippetHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("snippet_id", q.Get(":id"))
	r.URL.RawQuery = q.Encode()
	s.handle.RemoveSnippet(w, r)
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
		//Addr:     fmt.Sprintf(":%s", s.addr),
		Addr: fmt.Sprintf(":%s", s.addr),
		//ErrorLog: s.logger.ErrLog,
		ErrorLog: s.logger.ErrLog,
		Handler:  s.routes(),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	//s.logger.Info(fmt.Sprintf("server running on port:%s", s.addr))
	s.logger.Info(fmt.Sprintf("server running on port:%s", s.addr))

	// Call the ListenAndServe() method on our new http.Server struct.
	//err := srv.ListenAndServe()

	// Use the ListenAndServeTLS() method to start the HTTPS server. We
	// pass in the paths to the TLS certificate and corresponding private key as
	// the two parameters.
	err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

	//log error from serve
	//s.logger.Error(err.Error(), err)
	s.logger.Error(err.Error(), err)

	return err
}
