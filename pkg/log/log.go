package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var mtx sync.Mutex

type logLevel int

const (
	LDebug = logLevel(0)
	LInfo  = logLevel(1)
	LWarn  = logLevel(2)
	LErr   = logLevel(3)
	LFatal = logLevel(4)
)

func (l logLevel) Int() int {
	return int(l)
}

func (l logLevel) Name() string {

	switch l {
	case LDebug:
		return "DBG"
	case LInfo:
		return "NFO"
	case LWarn:
		return "WRN"
	case LErr:
		return "ERR"
	case LFatal:
		return "FTL"
	default:
		return "???"
	}

}

var level = LInfo

// SetLevel sets global logger level
func SetLevel(l logLevel) {
	if l >= LDebug && l <= LFatal {
		level = l
	}
}

var hook func(int, string)

func SetHook(h func(int, string)) {
	hook = h
}

var out io.Writer = os.Stdout

// SetOutput sets logging to output to a writer
func SetOutput(w io.Writer) {
	out = w
}

func accept(lvl logLevel) bool {
	return lvl.Int() >= level.Int()
}

func log(lvl logLevel, a ...interface{}) {
	if accept(lvl) {
		logMsg(lvl, fmt.Sprint(a...))
		if lvl == LFatal {
			panic(a)
		}
		if hook != nil {
			hook(lvl.Int(), fmt.Sprint(a...))
		}
	}
}

func logf(lvl logLevel, f string, args ...interface{}) {
	if accept(lvl) {
		m := fmt.Sprintf(f, args...)
		logMsg(lvl, m)
		if lvl == LFatal {
			panic(m)
		}
		if hook != nil {
			hook(lvl.Int(), m)
		}
	}
}

func logMsg(lvl logLevel, m string) {
	t := time.Now()
	mtx.Lock()
	fmt.Fprintf(out, "%s[%s] %s\t%s%s", lvl.Color(), lvl.Name(), t.Format(time.RFC3339), m, colorReset)
	fmt.Fprintln(out)
	mtx.Unlock()
}

// Debug prints debug message
func Debug(a ...interface{}) {
	log(LDebug, a...)
}

// Debugf prints formatted debug message
func Debugf(f string, args ...interface{}) {
	logf(LDebug, f, args...)
}

// Info prints info log message
func Info(a ...interface{}) {
	log(LInfo, a...)
}

// Infof prints formatted info log message
func Infof(f string, args ...interface{}) {
	logf(LInfo, f, args...)
}

// Warn prints warn log message
func Warn(a ...interface{}) {
	log(LWarn, a...)
}

// Warnf prints formatted warn log message
func Warnf(f string, args ...interface{}) {
	logf(LWarn, f, args...)
}

// Error prints error log message
func Error(a ...interface{}) {
	log(LErr, a...)
}

// Errorf prints formatted error log message
func Errorf(f string, args ...interface{}) {
	logf(LErr, f, args...)
}

// Fatal print fatal error and panic
func Fatal(a ...interface{}) {
	log(LFatal, a...)
}

// Fatalf print formatted fatal error and panic
func Fatalf(f string, args ...interface{}) {
	logf(LFatal, f, args...)
}
