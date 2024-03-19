package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"snippetbox/cmd/web/cache"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/forms"
	"snippetbox/pkg/logger"
	"strconv"
)

var (
	ErrInternalServerErr = errors.New("internal Server Error")
)

type Handle struct {
	logger        logger.Logger
	snippets      models.ISnippet
	templateCache map[string]*template.Template
}

func NewHandle(snippet models.ISnippet, lg logger.Logger) (*Handle, error) {

	// Initialize a new template cache...
	dir := filepath.Join(".", "ui", "html") // "./ui/html/"
	templateCache, err := cache.NewTemplateCache(dir)
	if err != nil {
		return nil, err
	}

	return &Handle{
		snippets:      snippet,
		logger:        lg,
		templateCache: templateCache,
	}, nil
}

func (h *Handle) Home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	//panic("oops! something went wrong") // Deliberate panic for testing

	s, err := h.snippets.Latest()
	if err != nil {
		h.serverError(w, err)
		return
	}

	// Use the new render helper.
	h.render(w, r, "home.page.tmpl", &cache.TemplateData{Snippets: s})

}
func (h *Handle) ShowSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("snippet_id"))
	if err != nil || id < 1 {
		h.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	s, err := h.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.notFound(w)
		} else {
			h.serverError(w, err)
		}
		return
	}
	// Use the new render helper.
	h.render(w, r, "show.page.tmpl", &cache.TemplateData{
		Snippet: s,
	})

}

func (h *Handle) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		h.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
		return
	}

	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError helper to send
	// a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.NewForm(r.PostForm)
	form.Required("title", "content", "expires")
	form.IsString("title", "content")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")
	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		h.render(w, r, "create.page.tmpl", &cache.TemplateData{Form: form})
		return
	}
	// Because the form data (with type url.Values) has been anonymously embedded
	// in the form.Form struct, we can use the Get() method to retrieve
	// the validated value for a particular form field.
	id, err := h.snippets.Insert(form.Values.Get("title"), form.Values.Get("content"), form.Values.Get("expires"))
	if err != nil {
		h.serverError(w, err)
		return
	}

	// Change the redirect to use the new semantic URL style of /snippet/:id
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

	//w.Write([]byte(fmt.Sprintf("Create a new snippet with id xxx")))
}
func (h *Handle) HealthChecker(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("health check ok"))
}

func (h *Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (h *Handle) CreateSnippetForm(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "create.page.tmpl", nil)
}
