package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"snippetbox/pkg/services"
)

type ILogger interface {
	Error(format string, args ...any)
	Info(format string, args ...any)
	Debug(format string, depth int)
	Fatal(format string, args ...any)
}

type StdLogger struct {
	ErrLog, InfoLog *log.Logger
}

func NewStdLogger() *StdLogger {

	//log info
	info := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime|log.LUTC)
	//log error
	err := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)

	return &StdLogger{
		ErrLog:  err,
		InfoLog: info,
	}
}

func (l StdLogger) Error(format string, args ...any) {
	l.ErrLog.Printf(format, args...)
}

func (l StdLogger) Info(format string, args ...any) {
	l.InfoLog.Printf(format, args...)
}

func (l StdLogger) Debug(format string, depth int) {
	trace := fmt.Sprintf("%s\n%s", format, debug.Stack())
	l.ErrLog.Output(depth, trace)
}

func (l StdLogger) Fatal(format string, args ...any) {
	l.ErrLog.Fatalf(format, args...)
}

//************* LOGGER FACTORY *********************///

var (
	ErrUnsupportedEnv    = errors.New("unsupported environment")
	ErrUnsupportedLogger = errors.New("unsupported logger")
)

const (
	LogInstanceStdLogger int = iota
	LogInstanceSlogLogger
)
const (
	EnvInstanceDev int = iota
	EnvInstanceProd
)

func NewLoggerFactory(envInstance, loggerInstance int, errLogFile, infoLogFile string) (ILogger, error) {

	switch envInstance {
	case EnvInstanceDev:
		switch loggerInstance {
		case LogInstanceSlogLogger:
			return nil, ErrUnsupportedLogger
		case LogInstanceStdLogger:
			//log info
			//infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime|log.LUTC)
			//log error
			//errLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)

			return NewStdLogger(), nil
		default:
			return nil, ErrUnsupportedLogger
		}
	case EnvInstanceProd:
		switch loggerInstance {
		case LogInstanceSlogLogger:
			return nil, ErrUnsupportedLogger
		case LogInstanceStdLogger:

			// get infoLog file or return nil and error if any
			errFile, err := fileWrite(errLogFile)
			//defer errFile.Close()
			if err != nil {
				return nil, err
			}

			//get errLog file or return error
			infoFile, err := fileWrite(infoLogFile)
			//defer infoFile.Close()
			if err != nil {
				return nil, err
			}

			//log info
			//infoLog := log.New(infoFile, "[INFO]\t", log.Ldate|log.Ltime|log.LUTC)
			//log error
			//errLog := log.New(errFile, "[ERROR]\t", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)

			os.Stdout = infoFile
			os.Stderr = errFile

			//return newLogger
			return NewStdLogger(), nil
		default:
			return StdLogger{}, ErrUnsupportedLogger
		}
	default:
		return StdLogger{}, ErrUnsupportedEnv
	}

}

func fileWrite(logFile string) (*os.File, error) {

	// Create any directories needed to put this file in them
	dirPath := "./logs/"
	dirFileMode := os.ModeDir | (services.OS_USER_RWX | services.OS_ALL_R)
	err := os.MkdirAll(dirPath, dirFileMode)
	if err != nil {
		return nil, err
	}

	//open file to log to and create if not existed
	file := filepath.Join(dirPath, logFile)
	fw, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return fw, err

}
