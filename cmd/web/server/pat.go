package server

import (
	"crypto/tls"
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/golangcollege/sessions"
	"github.com/justinas/alice"
	"net/http"
	"path/filepath"
	"snippetbox/cmd/web/handlers"
	"snippetbox/pkg/domain/ports"
	"snippetbox/pkg/logger"
	"time"
)

type Pat struct {
	router  *pat.PatternServeMux
	handle  *handlers.Handle
	logger  logger.ILogger
	addr    string
	session *sessions.Session
}

func NewPat(
	lg *logger.ILogger,
	addr string,
	user *ports.IUserRepository,
	snippet *ports.ISnippetRepository,
	session *sessions.Session,
	staticFileDir string,
) (*Pat, error) {

	//func NewPat(app *app.App) (*Pat, error) {

	//h, err := handlers.NewHandle(snippet, lg)
	h, err := handlers.NewHandle(user, snippet, lg, session, staticFileDir)
	if err != nil {
		return nil, err
	}

	// return server
	return &Pat{
		router:  pat.New(),
		logger:  *lg,
		addr:    addr,
		handle:  h,
		session: session,
	}, nil

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
	dynamicMiddleware := alice.New(s.session.Enable, s.handle.NoSurf, s.handle.Authenticate)

	// Create a file server which serves files out of the "./ui/static" directory.
	//	// Note that the path given to the http.Dir function is relative to the project
	//	// directory root.

	dir := filepath.Join(".", "ui", "static")
	fileServer := http.FileServer(http.Dir(dir))
	s.router.Get("/static", http.NotFoundHandler())
	s.router.Get("/static/", http.StripPrefix("/static", fileServer))

	//s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	s.router.Get("/health", http.HandlerFunc(s.handle.HealthChecker))

	//auth routes
	s.router.Get("/user/signup", dynamicMiddleware.ThenFunc(s.handle.SignupUserForm))
	s.router.Post("/user/signup", dynamicMiddleware.ThenFunc(s.handle.SignupUser))
	s.router.Get("/user/login", dynamicMiddleware.ThenFunc(s.handle.LoginUserForm))
	s.router.Post("/user/login", dynamicMiddleware.ThenFunc(s.handle.LoginUser))
	s.router.Post("/user/logout", dynamicMiddleware.Append(s.handle.RequireAuthentication).ThenFunc(s.handle.LogoutUser))

	// other routes
	s.router.Get("/", dynamicMiddleware.ThenFunc(s.handle.Home))
	s.router.Get("/snippet/create", dynamicMiddleware.Append(s.handle.RequireAuthentication).ThenFunc(s.handle.CreateSnippetForm))
	s.router.Post("/snippet/create", dynamicMiddleware.Append(s.handle.RequireAuthentication).ThenFunc(s.handle.CreateSnippet))
	s.router.Post("/snippet/remove/:id", dynamicMiddleware.Append(s.handle.RequireAuthentication).ThenFunc(s.removeSnippetHandleBuilder))
	s.router.Get("/snippet/:id", dynamicMiddleware.ThenFunc(s.showSnippetHandleBuilder))

	// Return the 'standard' middleware chain followed by the servemux router.
	return standardMiddleware.Then(s.router)
}

func (s *Pat) showSnippetHandleBuilder(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("snippet_id", q.Get(":id"))
	r.URL.RawQuery = q.Encode()
	s.handle.ShowSnippet(w, r)
}

func (s *Pat) removeSnippetHandleBuilder(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("snippet_id", q.Get(":id"))
	r.URL.RawQuery = q.Encode()
	s.handle.RemoveSnippet(w, r)
}

func (s *Pat) Begin() error {

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		//Addr:     fmt.Sprintf(":%s", s.addr),
		Addr: fmt.Sprintf(":%s", s.addr),
		//ErrorLog:  s.logger.ErrLog,
		Handler:   s.routes(),
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
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
