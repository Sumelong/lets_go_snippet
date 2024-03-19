package server

import (
	"errors"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
)

const (
	ServeInstancePat int = iota
	ServeInstanceGoMux
	ServeInstanceGorillaMux
	ServeInstanceGorillaPat
)

type IServer interface {
	Begin() error
}

var ErrUnsupportedServer = errors.New("unsupported server")

func NewServerFactory(serverInstance int, lg logger.Logger, addr string, snippet models.ISnippet) (IServer, error) {

	switch serverInstance {
	case ServeInstancePat:
		return NewPat(lg, addr, snippet)
	case ServeInstanceGoMux:
		return NewGoMux(lg, addr, snippet)
	case ServeInstanceGorillaMux:
		return NewGorillaMux(lg, addr, snippet)
	case ServeInstanceGorillaPat:
		return NewGorillaPat(lg, addr, snippet)
	default:
		return nil, ErrUnsupportedServer
	}
}
