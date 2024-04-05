package server

import (
	"bytes"
	"net/http"
	"net/url"
	"snippetbox/pkg/services"
	"testing"
)

func TestPat_routes(t *testing.T) {

	app := newTestApp(t)
	//staticFileDir := filepath.Join(".", "..", "..", "..", "ui", "html")
	staticFileDir, err := services.FindFile("ui/html")
	if err != nil {
		t.Fatalf("static files not found %s", err)
	}
	r, _ := NewPat(&app.logger, "", app.user, app.snippet, app.session, staticFileDir)
	ts := newTestServer(t, r.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		wantCode int
		wantBody []byte
		path     string
	}{
		{
			name:     "health_check",
			path:     "/health",
			wantCode: http.StatusOK,
			wantBody: []byte("health check ok"),
		},
		{
			name:     "snippet_create",
			path:     "/snippet/create",
			wantCode: http.StatusSeeOther,
		},
		{
			name:     "Valid ID",
			path:     "/snippet/1",
			wantCode: http.StatusOK,
			wantBody: []byte("An old silent pond..."),
		},
		{
			name:     "Non-existent ID",
			path:     "/snippet/2",
			wantCode: http.StatusNotFound,
			wantBody: nil,
		},
		{
			name:     "Negative ID",
			path:     "/snippet/-1",
			wantCode: http.StatusNotFound,
			wantBody: nil,
		},
		{
			name:     "Decimal ID",
			path:     "/snippet/1.23",
			wantCode: http.StatusNotFound,
			wantBody: nil,
		},
		{
			name:     "String ID",
			path:     "/snippet/foo",
			wantCode: http.StatusNotFound,
			wantBody: nil,
		},
		{
			name:     "Empty ID",
			path:     "/snippet/",
			wantCode: http.StatusNotFound,
			wantBody: nil,
		},
		{
			name:     "Trailing slash",
			path:     "/snippet/1/",
			wantCode: http.StatusNotFound,
			wantBody: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code, _, body := ts.get(t, tt.path)

			// We can then check the value of the response status code and body using
			// the same code as before.
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q; got %q", tt.wantBody, body)
			}

		})
	}
}

func TestPat_signUp(t *testing.T) {
	// Create the application struct containing our mocked dependencies and set
	// up the test server for running and end-to-end test.
	app := newTestApp(t)
	//staticFileDir := filepath.Join(".", "..", "..", "..", "ui", "html")
	staticFileDir, err := services.FindFile("ui/html")
	if err != nil {
		t.Fatalf("static files not found %s", err)
	}
	r, _ := NewPat(&app.logger, "", app.user, app.snippet, app.session, staticFileDir)
	ts := newTestServer(t, r.routes())
	defer ts.Close()

	// Make a GET /user/signup request and then extract the CSRF token from the
	// response body.
	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)
	// Log the CSRF token value in our test output. To see the output from the
	// t.Log() command you need to run `go test` with the -v (verbose) flag
	// enabled.
	t.Log(csrfToken)
}

func TestPat_SignupUser(t *testing.T) {
	app := newTestApp(t)
	staticFileDir, err := services.FindFile("ui/html")
	if err != nil {
		t.Fatalf("static files not found %s", err)
	}
	r, _ := NewPat(&app.logger, "", app.user, app.snippet, app.session, staticFileDir)
	ts := newTestServer(t, r.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)
	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantBody     []byte
	}{
		{
			"Valid submission",
			"Bob",
			"bob@example.com",
			"validPa$$word",
			csrfToken,
			http.StatusSeeOther,
			nil,
		},
		{
			"Empty name",
			"",
			"bob@example.com",
			"validPa$$word",
			csrfToken,
			http.StatusOK,
			[]byte("name cannot be empty"),
		},
		{
			"Empty email",
			"Bob",
			"",
			"validPa$$word",
			csrfToken,
			http.StatusOK,
			[]byte("email cannot be empty"),
		},
		{
			"Empty password",
			"Bob",
			"bob@example.com",
			"",
			csrfToken,
			http.StatusOK,
			[]byte("password cannot be empty"),
		},
		{
			"Invalid email (incomplete domain)",
			"Bob",
			"bob@example.",
			"validPa$$word",
			csrfToken,
			http.StatusOK,
			[]byte("This field is invalid"),
		},
		{
			"Invalid email (missing @)",
			"Bob",
			"bobexample.com",
			"validPa$$word",
			csrfToken,
			http.StatusOK,
			[]byte("This field is invalid"),
		},
		{
			"Invalid email (missing local part)",
			"Bob",
			"@example.com",
			"validPa$$word",
			csrfToken,
			http.StatusOK,
			[]byte("This field is invalid"),
		},
		{
			"Short password",
			"Bob",
			"bob@example.com",
			"pa$$word",
			csrfToken,
			http.StatusOK,
			[]byte("This field is too short (minimum is 10 characters)"),
		},
		{
			"Duplicate email",
			"Bob",
			"dupe@example.com",
			"validPa$$word",
			csrfToken,
			http.StatusOK,
			[]byte("Address is already in use"),
		},
		{
			"Invalid CSRF Token",
			"",
			"",
			"",
			"wrongToken",
			http.StatusBadRequest,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)
			code, _, body := ts.postForm(t, "/user/signup", form)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}
