package handlers

import (
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
	"snippetbox/cmd/web/cache"
	"time"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (h *Handle) serverError(w http.ResponseWriter, err error) {
	h.logger.Debug(err.Error(), 2)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Return true if the current request is from authenticated user, otherwise return false.
func (h *Handle) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(h.contextKeyIsAuthenticated).(bool)
	h.logger.Info("isAuth Value :", isAuthenticated)
	h.logger.Info("okay Value :", ok)
	if !ok {
		return false
	}
	return isAuthenticated
	//return true //h.session.Exists(r, "authenticatedUserID")

}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (h *Handle) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (h *Handle) notFound(w http.ResponseWriter) {
	h.clientError(w, http.StatusNotFound)
}

// Create an addDefaultData helper. This takes a pointer to a templateData
// struct, adds the current year to the CurrentYear field, and then returns
// the pointer. Again, we're not using the *http.Request parameter at the
// moment, but we will do later in the book.
func (h *Handle) addDefaultData(td *cache.TemplateData, r *http.Request) *cache.TemplateData {
	if td == nil {
		td = &cache.TemplateData{}
	}

	// Add the CSRF token to the templateData struct.
	td.CSRFToken = nosurf.Token(r)

	td.CurrentYear = time.Now().Year()

	// Use the PopString() method to retrieve the value for the "flash" key.
	// PopString() also deletes the key and value from the session data, so it
	// acts like a one-time fetch. If there is no matching key in the session
	// data this will return the empty string.
	//
	// Add the flash message to the template data, if one exists.
	td.Flash = h.session.PopString(r, "flash")

	// Add the authentication status to the template data.
	td.IsAuthenticated = h.isAuthenticated(r)

	return td
}

func (h *Handle) render(w http.ResponseWriter, r *http.Request, name string, td *cache.TemplateData) {

	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the
	// provided name, call serverError.
	ts, ok := h.templateCache[name]
	if !ok {
		h.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)
	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError helper and then
	// return.
	err := ts.Execute(buf, h.addDefaultData(td, r))
	if err != nil {
		h.serverError(w, err)
		return
	}
	// Write the contents of the buffer to the http.ResponseWriter. Again, this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	buf.WriteTo(w)

}
