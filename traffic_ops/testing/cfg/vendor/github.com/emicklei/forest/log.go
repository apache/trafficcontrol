package forest

import (
	"runtime"
	"strings"
)

var scanStackForFile = true
var logf_func = logf

const noStackOffset = 0

// Logf adds the actual file:line information to the log message
func Logf(t T, format string, args ...interface{}) {
	logf_func(t, noStackOffset, "\n"+format, args...)
}

func logfatal(t T, format string, args ...interface{}) {
	logf_func(t, noStackOffset, format, args...)
	t.FailNow()
}

func logerror(t T, format string, args ...interface{}) {
	logf_func(t, noStackOffset, format, args...)
	t.Fail()
}

func logf(t T, stackOffset int, format string, args ...interface{}) {
	var file string
	var line int
	var ok bool
	if scanStackForFile {
		offset := 0
		outside := false
		for !outside {
			_, file, line, ok = runtime.Caller(2 + offset)
			outside = !strings.Contains(file, "/forest/")
			offset++
		}
	} else {
		_, file, line, ok = runtime.Caller(2)
	}
	if ok {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}
	t.Logf("<-- %s:%d"+format, append([]interface{}{file, line}, args...)...)
}
