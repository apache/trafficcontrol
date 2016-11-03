// Inspired by https://www.goinggo.net/2013/11/using-log-package-in-go.html
package log

import (
	"fmt"
	"io"
	"log"
	"time"
)

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(errW, warnW, infoW, debugW io.Writer) {
	Debug = log.New(debugW, "DEBUG: ", log.Lshortfile)
	Info = log.New(infoW, "INFO: ", log.Lshortfile)
	Warning = log.New(warnW, "WARNING: ", log.Lshortfile)
	Error = log.New(errW, "ERROR: ", log.Lshortfile)
}

const timeFormat = time.RFC3339Nano

func Errorf(format string, v ...interface{}) {
	Error.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintf(format, v...))
}
func Errorln(v ...interface{}) {
	Error.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintln(v...))
}
func Warnf(format string, v ...interface{}) {
	Warning.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintf(format, v...))
}
func Warnln(v ...interface{}) {
	Warning.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintln(v...))
}
func Infof(format string, v ...interface{}) {
	Info.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintf(format, v...))
}
func Infoln(v ...interface{}) {
	Info.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintln(v...))
}
func Debugf(format string, v ...interface{}) {
	Debug.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintf(format, v...))
}
func Debugln(v ...interface{}) {
	Debug.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintln(v...))
}

// Close calls `Close()` on the given Closer, and logs any error. On error, the context is logged, followed by a colon, the error message, and a newline. This is primarily designed to be used in `defer`, for example, `defer log.Close(resp.Body, "readData fetching /foo/bar")`.
func Close(c io.Closer, context string) {
	err := c.Close()
	if err != nil {
		Errorf("%v: %v", context, err)
	}
}

// Closef acts like Close, with a given format string and values, followed by a colon, the error message, and a newline. The given values are not coerced, concatenated, or printed unless an error occurs, so this is more efficient than `Close()`.
func Closef(c io.Closer, contextFormat string, v ...interface{}) {
	err := c.Close()
	if err != nil {
		Errorf(contextFormat, v...)
		Errorf(": %v", err)
	}
}
