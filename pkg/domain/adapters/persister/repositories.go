package persister

/*
// Repositories Define a map to hold repositories for different types
type Repositories struct {
	repositories map[string]ports.Repository[any]
	lg           logger.ILogger
}

// NewRepositories creates a new instance of Repositories
func NewRepositories(lg *logger.ILogger) *Repositories {
	return &Repositories{
		repositories: make(map[string]ports.Repository[any]),
		lg:           *lg,
	}
}

// RegisterRepository adds a new repository for a specific type
func (rm *Repositories) RegisterRepository(name string, r ports.Repository[any]) {

	repo, ok := r.(ports.Repository[any])
	if !ok {
		rm.lg.Error("failed to convert %s", r)
		return
	}
	rm.repositories[name] = repo

}

// GetRepository retrieves a registered repository by name
func (rm *Repositories) GetRepository(name string) (ports.Repository[any], error) {
	r, ok := rm.repositories[name]
	if !ok {
		return nil, errors.New("repository not found")
	}
	return r, nil
}
*/
