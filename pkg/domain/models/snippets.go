package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("domain: no matching record found")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type ISnippet interface {
	Insert(title, content, expires string) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}
