package logger

import (
	"io"
	"log"
)

type MockLogger struct {
	logger  ILogger
	ErrLog  *log.Logger
	InfoLog *log.Logger
}

func NewMockLogger() ILogger {
	return MockLogger{
		logger:  nil,
		ErrLog:  log.New(io.Discard, "", 0),
		InfoLog: log.New(io.Discard, "", 0),
	}
}

func (m MockLogger) Error(format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (m MockLogger) Info(format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (m MockLogger) Debug(format string, depth int) {
	//TODO implement me
	panic("implement me")
}

func (m MockLogger) Fatal(format string, args ...any) {
	//TODO implement me
	panic("implement me")
}
