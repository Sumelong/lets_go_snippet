package handlers

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"snippetbox/cmd/web/cache"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/forms"
)

func (h *Handle) SignupUserForm(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "signup.page.tmpl", &cache.TemplateData{
		Form: forms.NewForm(nil),
	})

}
func (h *Handle) SignupUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data.
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}
	// Validate the form contents using the form helper we made earlier.
	form := forms.NewForm(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)
	// If there are any errors, redisplay the signup form.
	if !form.Valid() {
		h.render(w, r, "signup.page.tmpl", &cache.TemplateData{Form: form})
		return
	}

	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Values.Get("password")), 12)

	newUser := models.User{
		Name:           form.Values.Get("name"),
		Email:          form.Values.Get("email"),
		HashedPassword: hashedPassword,
	}

	// Try to create a new user record in the database. If the email already exists
	// add an error message to the form and re-display it.
	_, err = h.user.Create(newUser)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			h.render(w, r, "signup.page.tmpl", &cache.TemplateData{Form: form})
		} else {
			h.serverError(w, err)
		}
		return
	}
	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked and asking them to log in.
	h.session.Put(r, "flash", "Your signup was successful. Please log in.")
	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (h *Handle) LoginUserForm(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "login.page.tmpl", &cache.TemplateData{
		Form: forms.NewForm(nil),
	})
}
func (h *Handle) LoginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic error
	// message to the form failures map and re-display the login page.
	form := forms.NewForm(r.PostForm)
	id, err := h.user.Authenticate(form.Values.Get("email"), form.Values.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			h.render(w, r, "login.page.tmpl", &cache.TemplateData{Form: form})
		} else {
			h.serverError(w, err)
		}
		return
	}
	// Add the ID of the current user to the session, so that they are now 'logged
	// in'.
	h.session.Put(r, "authenticatedUserID", id)
	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}
func (h *Handle) LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	h.session.Remove(r, "authenticatedUserID")
	// Add a flash message to the session to confirm to the user that they've been
	// logged out.
	h.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
