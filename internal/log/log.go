package log

import (
	"github.com/withmandala/go-log"
	"os"
)

type ILogger interface {
	Init()
	Info(message ...interface{})
	Warn(message ...interface{})
	Error(message ...interface{})
	Debug(message ...interface{})
	Trace(message ...interface{})
	Fatal(message ...interface{})
}

func InitLogger() *log.Logger {
	return log.New(os.Stderr).WithTimestamp().WithColor().WithDebug()
}
