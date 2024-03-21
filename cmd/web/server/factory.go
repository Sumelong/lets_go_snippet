package server

import (
	"errors"
	"github.com/golangcollege/sessions"
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

func NewServerFactory(
	serverInstance int,
	lg *logger.Logger,
	addr string,
	snippet *models.ISnippet,
	session *sessions.Session,
) (IServer, error) {

	switch serverInstance {
	case ServeInstancePat:
		return NewPat(lg, addr, snippet, session)
	case ServeInstanceGoMux:
		return NewGoMux(lg, addr, snippet, session)
	case ServeInstanceGorillaMux:
		return NewGorillaMux(lg, addr, snippet, session)
	case ServeInstanceGorillaPat:
		return NewGorillaPat(lg, addr, snippet, session)
	default:
		return nil, ErrUnsupportedServer
	}
}

/*
func NewServerFactory(serverInstance int, app *app.App) (IServer, error) {

	switch serverInstance {
	case ServeInstancePat:
		return NewPat(app)
	default:
		return nil, ErrUnsupportedServer
	}
}
*/
