package server

import (
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
	"snippetbox/cmd/web/controller"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
	"time"
)

type Pat struct {
	router     *pat.PatternServeMux
	controller *controller.Controller
	logger     logger.Logger
	addr       string
}

func NewPat(lg logger.Logger, addr string, snippet models.ISnippet) (*Pat, error) {

	c, err := controller.NewController(snippet, lg)
	if err != nil {
		return nil, err
	}

	// return server
	return &Pat{
		router:     pat.New(),
		logger:     lg,
		addr:       addr,
		controller: c,
	}, nil
}

func (s *Pat) routes() http.Handler {
	//flag.StringVar(&m.app.port, "addr", "4000", "HTTP network address")
	//flag.Parse()

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(s.controller.RecoverPanic, s.controller.LogRequest, s.controller.SecureHeaders)

	// Create a file server which serves files out of the "./ui/static" directory.
	//	// Note that the path given to the http.Dir function is relative to the project
	//	// directory root.

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	s.router.Get("/static", http.NotFoundHandler())
	s.router.Get("/static/", http.StripPrefix("/static", fileServer))

	s.router.Get("/health", http.HandlerFunc(s.controller.HealthChecker))
	s.router.Get("/", http.HandlerFunc(s.controller.Home))

	//So to ensure that the exact match takes preference,
	//we register the exact match routes before any wildcard routes.
	s.router.Get("/snippet/create", http.HandlerFunc(s.controller.CreateSnippetForm))
	s.router.Post("/snippet/create", http.HandlerFunc(s.controller.CreateSnippet))

	//So to ensure that the exact match takes preference,
	//we register the wildcard routes.
	s.router.Get("/snippet/:id", http.HandlerFunc(s.controller.ShowSnippet))

	// Return the 'standard' middleware chain followed by the servemux router.
	return standardMiddleware.Then(s.router)

}

func (s *Pat) Begin() error {

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
