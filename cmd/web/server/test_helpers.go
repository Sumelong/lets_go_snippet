package server

import (
	"database/sql"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/golangcollege/sessions"
	_ "modernc.org/sqlite"
	"snippetbox/pkg/domain/ports"
	"snippetbox/pkg/domain/ports/mocks"
	"snippetbox/pkg/logger"
	"snippetbox/pkg/services"
)

type testApp struct {
	logger        logger.ILogger
	user          *ports.IUserRepository
	snippet       *ports.ISnippetRepository
	session       *sessions.Session
	templateCache map[string]*template.Template
}

// Create a newTestApplication helper which returns an instance of our
// application struct containing mocked dependencies.
func newTestApp(t *testing.T) *testApp {
	//arrange controller dependencies
	var lg logger.ILogger = logger.StdLogger{
		ErrLog:  log.New(io.Discard, "", 0),
		InfoLog: log.New(io.Discard, "", 0),
	}
	var ur ports.IUserRepository = mocks.MockUserRepository{}
	var sr ports.ISnippetRepository = mocks.MockSnippetRepository{}

	// Create an instance of the template cache.
	/*templateCache, err := cache.NewTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}*/
	// Create a session manager instance, with the same settings as production.
	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	return &testApp{
		logger:  lg,
		user:    &ur,
		snippet: &sr,
		session: session,
		//templateCache: templateCache,
	}
}

// Define a custom TestServer type which anonymously embeds a httptest.Server
// instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initalizes and returns a new instance
// of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the cookie jar to the client, so that response cookies are stored
	// and then sent with subsequent requests.
	ts.Client().Jar = jar

	// Disable redirect-following for the client. Essentially this function
	// is called after a 3xx response is received by the client, and returning
	// the http.ErrUseLastResponse error forces it to immediately return the
	// received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

// Implement a Get method on our custom TestServer type. This makes a GET
// request to a given url path on the test server, and returns the response
// status code, headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) (statusCode int, header http.Header, body []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err = io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body
}

// Define a regular expression which captures the CSRF token value from the
// HTML for our user signup page.
var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

func extractCSRFToken(t *testing.T, body []byte) string {
	// Use the FindSubmatch method to extract the token from the HTML body.
	// Note that this returns an array with the entire matched pattern in the
	// first position, and the values of any captured data in the subsequent
	// positions.
	matches := csrfTokenRX.FindSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
	return html.UnescapeString(string(matches[1]))
}

// Create a postForm method for sending POST requests to the test server.
// The final parameter to this method is a url.Values object which can contain
// any data that you want to send in the request body.
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, []byte) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}
	// Read the response body.
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	// Return the response status, headers and body.
	return rs.StatusCode, rs.Header, body
}

func newTestDB(t *testing.T) (*sql.DB, func()) {
	// Establish a sql.DB connection pool for our test database. Because our
	// setup and teardown scripts contains multiple SQL statements, we need
	// to use the `multiStatements=true` parameter in our DSN. This instructs
	// our MySQL database driver to support executing multiple SQL statements
	// in one db.Exec()` call.
	// Connect to the SQLite database
	db, err := sql.Open("sqlite", "test_snippetbox.db")
	if err != nil {
		t.Fatalf("Error connecting to the database: %v\n", err)

	}
	err = db.Ping()
	if err != nil {
		t.Fatalf("Error opening database: %v\n", err)
	}

	// Read the setup SQL script from file and execute the statements.
	setUpScript, err := services.FindFile("testdata/setup.sql")
	if err != nil {
		t.Fatal("could not locate setup.sql script")
	}
	//script, err := os.ReadFile("./testdata/teardown.sql")
	script, err := os.ReadFile(setUpScript)
	//script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}
	// Return the connection pool and an anonymous function which reads and
	// executes the teardown script, and closes the connection pool. We can
	// assign this anonymous function and call it later once our test has
	// completed.
	return db, func() {
		tearDownScript, err := services.FindFile("testdata/teardown.sql")
		if err != nil {
			t.Fatal("could not locate teardown.sql script")
		}
		//script, err := os.ReadFile("./testdata/teardown.sql")
		script, err := os.ReadFile(tearDownScript)
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
}
