package server

import (
	"errors"
	"github.com/golangcollege/sessions"
	"snippetbox/pkg/domain/ports"
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
	lg *logger.ILogger,
	addr string,
	user *ports.IUserRepository,
	snippet *ports.ISnippetRepository,
	session *sessions.Session,
	staticFileDir string,
) (IServer, error) {

	switch serverInstance {
	case ServeInstancePat:
		return NewPat(lg, addr, user, snippet, session, staticFileDir)
	case ServeInstanceGoMux:
		return NewGoMux(lg, addr, user, snippet, session, staticFileDir)
	case ServeInstanceGorillaMux:
		return NewGorillaMux(lg, addr, user, snippet, session, staticFileDir)
	case ServeInstanceGorillaPat:
		return NewGorillaPat(lg, addr, user, snippet, session, staticFileDir)
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
