package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/cmd/web/cache"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/forms"
	"strconv"
)

func (h *Handle) ShowSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("snippet_id"))
	if err != nil || id < 1 {
		h.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	s, err := h.snippet.ReadOne(id)
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

func (h *Handle) CreateSnippetForm(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "create.page.tmpl", &cache.TemplateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.NewForm(nil),
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
	form.Required("title", "content", "expiresIn")
	form.IsString("title")
	form.MaxLength("title", 100)
	form.PermittedValues("expiresIn", "365", "7", "1")
	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		h.render(w, r, "create.page.tmpl", &cache.TemplateData{Form: form})
		return
	}
	newSnippet := models.Snippet{
		ID:        0,
		Title:     form.Values.Get("title"),
		Content:   form.Values.Get("content"),
		ExpiresIn: form.Values.Get("expiresIn"),
	}

	// Because the form data (with type url.Values) has been anonymously embedded
	// in the form.Form struct, we can use the Get() method to retrieve
	// the validated value for a particular form field.
	id, err := h.snippet.Create(newSnippet)
	if err != nil {
		h.serverError(w, err)
		return
	}

	// Use the Put() method to add a string value ("Your snippet was saved
	// successfully!") and the corresponding key ("flash") to the session
	// data. Note that if there's no existing session for the current user
	// (or their session has expired) then a new, empty, session for them
	// will automatically be created by the session middleware.
	h.session.Put(r, "flash", "Snippet successfully created!")

	// Change the redirect to use the new semantic URL style of /snippet/:id
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

	//w.Write([]byte(fmt.Sprintf("Create a new snippet with id xxx")))
}

func (h *Handle) RemoveSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("snippet_id"))
	if err != nil || id < 1 {
		h.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	res, err := h.snippet.Delete(uint(id))
	if err != nil {
		h.notFound(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	// Use the new render helper.
	h.render(w, r, "show.page.tmpl", &cache.TemplateData{
		Snippet: &models.Snippet{ID: int(res)},
	})

}

func (h *Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
