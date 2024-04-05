package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
	"snippetbox/pkg/domain/models"
)

func (h *Handle) SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

func (h *Handle) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect them to the login page and
		// return from the middleware chain so that no subsequent handlers in
		// the chain are executed.
		if !h.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")
		// And call the next handler in the chain.
		next.ServeHTTP(w, r)

	})
}

func (h *Handle) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check if a authenticatedUserID value exists in the session. If this *isn't
		// present* then call the next handler in the chain as normal.
		exist := h.session.Exists(r, "authenticatedUserID")
		h.logger.Info("authenticatedUserID exist? %v", exist)
		if !exist {
			next.ServeHTTP(w, r)
			return
		}

		// Fetch the details of the current user from the database. If no matching
		// record is found, or the current user is has been deactivated, remove the
		// (invalid) authenticatedUserID value from their session and call the next
		// handler in the chain as normal.

		userId := h.session.GetInt(r, "authenticatedUserID")

		h.logger.Info("userid gotten from session authenticatedUserID is- %v", userId)
		user, err := h.user.ReadOne(userId)
		h.logger.Info("user gotten : %v", user)
		h.logger.Info("error gotten : %v", err)
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			h.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			h.serverError(w, err)
			return
		}

		// Otherwise, we know that the request is coming from a active, authenticated,
		// user. We create a new copy of the request, with a true boolean value
		// added to the request context to indicate this, and call the next handler
		// in the chain *using this new copy of the request*.
		h.contextKeyIsAuthenticated = "isAuthenticated"
		ctx := context.WithValue(r.Context(), h.contextKeyIsAuthenticated, true)

		h.logger.Info("contextkeyis: %v", r.Context().Value(h.contextKeyIsAuthenticated))
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// NoSurf Create a NoSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly flags set.
func (h *Handle) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func (h *Handle) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("%s - %s %s %s\n", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (h *Handle) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				h.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
