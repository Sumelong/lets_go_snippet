package models

import "time"

type Snippetx struct {
	BaseModel

	Title   string
	Content string
	Expires time.Time
}

func (u Snippet) GetID() uint {
	return uint(u.ID)
}
