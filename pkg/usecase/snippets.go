package usecase

import (
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
)

type snippets struct {
	logger   logger.ILogger
	snippets models.ISnippet
}

func (s snippets) GetLatest() ([]*models.Snippet, error) {
	data, err := s.snippets.Latest()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s snippets) GetOne(id string) {

}
