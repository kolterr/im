package log

import (
	"github.com/sirupsen/logrus"
)

var (
	Prefix = "[ IM ] "
)

var ins *logrus.Logger

func SetPrefix(prefix string) {
	Prefix = prefix
}

func SetLevel(level logrus.Level) {
	ins.SetLevel(level)
}

func init() {
	ins = logrus.New()
}

func Debug(args ...interface{}) {
	ins.Debug(Prefix, args)
}

func Debugf(format string, args ...interface{}) {
	ins.Debugf(Prefix+format, args)
}

func Error(args ...interface{}) {
	ins.Error(Prefix, args)
}
func Errorf(format string, args ...interface{}) {
	ins.Errorf(Prefix+format, args)
}

func Info(args ...interface{}) {
	ins.Info(Prefix, args)
}

func Infof(format string, args ...interface{}) {
	ins.Infof(Prefix+format, args)
}

func Panic(args ...interface{}) {
	ins.Panic(args)
}
