package pkg

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

const (
	LogInstanceStdLogger int = iota
	LogInstanceSlogLogger
)

var ErrUnsupportedLogger = errors.ErrUnsupported

type ILogger interface {
	Error(format string, args ...any)
	Info(format string, args ...any)
	Debug(format string, args ...any)
	Fatal(format string, args ...any)
}

type Logger struct {
	ErrLog, InfoLog *log.Logger
}

func NewLogger() Logger {

	//log info
	info := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime|log.LUTC)
	//log error
	err := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)

	return Logger{
		ErrLog:  err,
		InfoLog: info,
	}
}

func (l Logger) Error(format string, args ...any) {
	l.ErrLog.Printf(format, args)
}

func (l Logger) Info(format string, args ...any) {
	l.InfoLog.Printf(format, args)
}

func (l Logger) Debug(format string, args ...any) {
	l.ErrLog.Printf(format, args)
}

func (l Logger) Fatal(format string, args ...any) {
	l.ErrLog.Fatalf(format, args)
}

//************* LOGGER FACTORY *********************///

func NewLoggerFactory(app *App) (Logger, error) {

	switch app.envInstance {
	case EnvInstanceDev:
		switch app.loggerInstance {
		case LogInstanceSlogLogger:
			return Logger{}, ErrUnsupportedLogger
		case LogInstanceStdLogger:
			//log info
			//infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime|log.LUTC)
			//log error
			//errLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)

			return NewLogger(), nil
		default:
			return Logger{}, ErrUnsupportedLogger
		}
	case EnvInstanceProd:
		switch app.loggerInstance {
		case LogInstanceSlogLogger:
			return Logger{}, ErrUnsupportedLogger
		case LogInstanceStdLogger:

			// get infoLog file or return nil and error if any
			errFile, err := fileWrite(app.prodErrLogFile)
			//defer errFile.Close()
			if err != nil {
				return Logger{}, err
			}

			//get errLog file or return error
			infoFile, err := fileWrite(app.prodInfoLogFile)
			//defer infoFile.Close()
			if err != nil {
				return Logger{}, err
			}

			//log info
			//infoLog := log.New(infoFile, "[INFO]\t", log.Ldate|log.Ltime|log.LUTC)
			//log error
			//errLog := log.New(errFile, "[ERROR]\t", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)

			os.Stdout = infoFile
			os.Stderr = errFile

			//return newLogger
			return NewLogger(), nil
		default:
			return Logger{}, ErrUnsupportedLogger
		}
	default:
		return Logger{}, ErrUnsupportedEnv
	}

}

func fileWrite(logFile string) (*os.File, error) {

	// Create any directories needed to put this file in them
	dirPath := "./logs/"
	dirFileMode := os.ModeDir | (OS_USER_RWX | OS_ALL_R)
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
