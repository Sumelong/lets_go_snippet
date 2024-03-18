package controller

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
	"strconv"
)

var (
	ErrInternalServerErr = errors.New("internal Server Error")
)

type Controller struct {
	logger        logger.Logger
	snippets      models.ISnippet
	templateCache map[string]*template.Template
}

func NewController(snippet models.ISnippet, lg logger.Logger) (*Controller, error) {

	// Initialize a new template cache...
	dir := filepath.Join(".", "ui", "html") // "./ui/html/"
	templateCache, err := NewTemplateCache(dir)
	if err != nil {
		return nil, err
	}

	return &Controller{
		snippets:      snippet,
		logger:        lg,
		templateCache: templateCache,
	}, nil
}

func (c *Controller) Home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	//panic("oops! something went wrong") // Deliberate panic for testing

	s, err := c.snippets.Latest()
	if err != nil {
		c.serverError(w, err)
		return
	}

	// Use the new render helper.
	c.render(w, r, "home.page.tmpl", &templateData{Snippets: s})

}
func (c *Controller) ShowSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("snippet_id"))
	if err != nil || id < 1 {
		c.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	s, err := c.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			c.notFound(w)
		} else {
			c.serverError(w, err)
		}
		return
	}
	// Use the new render helper.
	c.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})

}
func (c *Controller) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		c.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
		//http.Error(w, "Method Not Allowed", 405)
		return
	}

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"
	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := c.snippets.Insert(title, content, expires)
	if err != nil {
		c.serverError(w, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
	//http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

	// Change the redirect to use the new semantic URL style of /snippet/:id
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

	//w.Write([]byte(fmt.Sprintf("Create a new snippet with id xxx")))
}
func (c *Controller) HealthChecker(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("health check ok"))
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) CreateSnippetForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet..."))

}
