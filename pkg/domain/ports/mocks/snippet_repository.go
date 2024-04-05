package mocks

import (
	"snippetbox/pkg/domain/models"
	"time"
)

var mockSnippet = &models.Snippet{
	ID:        1,
	Title:     "An old silent pond",
	Content:   "An old silent pond...",
	CreatedOn: time.Now(),
	ExpiresOn: time.Now(),
	ExpiresIn: "5",
}

type MockSnippetRepository struct{}

func (u MockSnippetRepository) Create(e models.Snippet) (uint, error) {
	return uint(e.ID), nil
}

func (u MockSnippetRepository) ReadAll() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}

func (u MockSnippetRepository) ReadOne(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (u MockSnippetRepository) ReadBy(e *models.Snippet) ([]*models.Snippet, error) {
	if e != nil {
		return []*models.Snippet{mockSnippet}, nil
	}
	return nil, models.ErrNoRecord
}

func (u MockSnippetRepository) Update(e *models.Snippet) (uint, error) {
	if e == nil {
		return 0, models.ErrNoRecord
	}
	mockSnippet.ID = e.ID
	mockSnippet.Title = e.Title
	mockSnippet.Content = e.Content
	mockSnippet.CreatedOn = e.CreatedOn
	mockSnippet.ExpiresOn = e.ExpiresOn
	mockSnippet.ExpiresIn = e.ExpiresIn
	return uint(mockSnippet.ID), nil
}

func (u MockSnippetRepository) Delete(id uint) (uint, error) {
	switch id {
	case 1:
		return id, nil
	default:
		return 0, models.ErrNoRecord
	}
}
