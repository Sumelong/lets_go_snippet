package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord             = errors.New("model: no matching record found")
	ErrInvalidCredentials   = errors.New("models: invalid credentials")
	ErrDuplicateEmail       = errors.New("models: duplicate email")
	ErrEntityCreationFailed = errors.New("models: failed to add entity")
)

type BaseModel struct {
	ID      int
	Created time.Time
	updated time.Time
}

type Snippet struct {
	ID        int
	Title     string
	Content   string
	CreatedOn time.Time
	ExpiresOn time.Time
	ExpiresIn string
}

type ISnippet interface {
	Insert(title, content, expires string) (int, error)
	Get(id int) (*Snippet, error)
	Remove(id int) (int, error)
	Latest() ([]*Snippet, error)
}
