package gmx

// pkg/os instrumentation

import (
	"os"
)

func init() {
	Publish("os.args", osArgs)
}

func osArgs() interface{} {
	return os.Args
}
