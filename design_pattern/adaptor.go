package designpattern

import (
	"log"
	"strings"
)

//
type Level int8

const (
	Debug Level = iota
	Info
	Warn
	Error
)

func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	default:
		return ""
	}
}

func parseLevel(level string) Level {
	s := strings.ToUpper(level)
	switch s {
	case "DEBUG":
		return Debug
	case "INFO":
		return Info
	case "WARN":
		return Warn
	case "ERROR":
		return Error
	default:
		return Info
	}
}

// Logger logger interface
type Logger interface {
	Log(level Level, kvs ...interface{}) error
}

type Zap struct {
}

func (z *Zap) Log(level Level, kvs ...interface{}) {
	log.Println("level1", level, kvs)
}

type Logrus struct {
}

func (z *Logrus) Log(level Level, kvs ...interface{}) {
	log.Println("level：", level, kvs)
}

type ZapAdaptor struct {
	log *Zap
}

func (z *ZapAdaptor) ZapAdaptor(level Level, kvs ...interface{}) {

}
