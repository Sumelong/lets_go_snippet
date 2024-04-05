package handlers

import (
	"errors"
	"github.com/golangcollege/sessions"
	"html/template"
	"net/http"
	"snippetbox/cmd/web/cache"
	"snippetbox/pkg/domain/ports"
	"snippetbox/pkg/logger"
)

var (
	ErrInternalServerErr = errors.New("internal Server Error")
)

type appContextKey string

//const contextKeyIsAuthenticated = appContextKey("isAuthenticated")

type Handle struct {
	logger                    logger.ILogger
	user                      ports.IUserRepository
	snippet                   ports.ISnippetRepository
	session                   *sessions.Session
	contextKeyIsAuthenticated appContextKey
	templateCache             map[string]*template.Template
}

func NewHandle(user *ports.IUserRepository, snippet *ports.ISnippetRepository,
	lg *logger.ILogger, session *sessions.Session, staticFileDir string) (*Handle, error) {

	// Initialize a new template cache...
	//dir := filepath.Join(".", "ui", "html") // "./ui/html/"
	dir := staticFileDir //filepath.Join(".", "ui", "html") // "./ui/html/"
	templateCache, err := cache.NewTemplateCache(dir)
	if err != nil {
		return nil, err
	}

	//contextKeyIsAuthenticated := appContextKey("isAuthenticated")

	return &Handle{
		snippet:       *snippet,
		user:          *user,
		logger:        *lg,
		session:       session,
		templateCache: templateCache,
		//contextKeyIsAuthenticated: "isAuthenticated",
	}, nil
}

func (h *Handle) HealthChecker(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("health check ok"))
}

func (h *Handle) Home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	//panic("oops! something went wrong") // Deliberate panic for testing

	//s, err := h.snippets.Latest()
	s, err := h.snippet.ReadAll()
	if err != nil {
		h.serverError(w, err)
		return
	}

	// Use the new render helper.
	h.render(w, r, "home.page.tmpl", &cache.TemplateData{Snippets: s})

}
