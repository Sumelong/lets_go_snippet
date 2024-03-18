package server

import (
	"errors"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
)

const (
	ServeInstancePat int = iota
	ServeInstanceMux
)

type IServer interface {
	Begin() error
}

var ErrUnsupportedServer = errors.New("unsupported server")

func NewServerFactory(serverInstance int, lg logger.Logger, addr string, snippet models.ISnippet) (IServer, error) {

	switch serverInstance {
	case ServeInstancePat:
		return NewPat(lg, addr, snippet)
	case ServeInstanceMux:
		return NewMux(lg, addr, snippet)
	default:
		return nil, ErrUnsupportedServer
	}
}
