package pkg

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var (
	ErrInternalServerErr = errors.New("internal Server Error")
)

type Handlers struct {
	lg Logger
}

func (h Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message and use
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		//log error
		//h.lg.Info(err.Error(), err)

		h.serverError(w, err) // Use the serverError() helper.

		//return error
		//http.Error(w, ErrInternalServerErr.Error(), http.StatusInternalServerError)
		return
	}
	// We then use the Execute() method on the template set to write the template
	// content as the response body. The last parameter to Execute() represents any
	// dynamic data that we want to pass in, which for now we'll leave as nil.
	err = ts.Execute(w, nil)
	if err != nil {
		h.serverError(w, err) // Use the serverError() helper.
	}

}
func (h Handlers) ShowSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		//return error
		http.NotFound(w, r)
		//log error
		h.lg.Info(err.Error(), err)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}
func (h Handlers) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		h.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
		//http.Error(w, "Method Not Allowed", 405)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
func (h Handlers) HealthChecker(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("health check ok"))
}

func (h Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
