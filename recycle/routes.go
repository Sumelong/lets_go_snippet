package recycle

/*
func routes(s *server.GoMux) http.Handler {
	//flag.StringVar(&m.app.port, "addr", "4000", "HTTP network address")
	//flag.Parse()

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(s.handlers.recoverPanic, s.handlers.logRequest, secureHeaders)
	r := pat.New()

	// Create a file server which serves files out of the "./ui/static" directory.
	//	// Note that the path given to the http.Dir function is relative to the project
	//	// directory root.
	fileServer := http.FileServer(neuteredFileSystem{fs: http.Dir("./ui/static/")})
	//fileServer := http.FileServer(http.Dir("./ui/static/"))
	s.router.Handle("/static", http.NotFoundHandler())
	s.router.Handle("/static/", http.StripPrefix("/static", fileServer))

	s.router.HandleFunc("/health", s.handlers.HealthChecker)
	s.router.HandleFunc("/", s.handlers.Home)
	s.router.HandleFunc("/snippet/ ", s.handlers.ShowSnippet)
	s.router.HandleFunc("/snippet/create", s.handlers.CreateSnippet)

	// Wrap the existing chain with the logRequest middleware.
	//return s.handlers.recoverPanic(s.handlers.logRequest(secureHeaders(s.router)))

	// Return the 'standard' middleware chain followed by the servemux router.
	return standardMiddleware.Then(s.router)

}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err = nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
*/
