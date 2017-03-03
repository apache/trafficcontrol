package nagios

import (
	"fmt"
	"os"
	"strings"
)

type Status int

const (
	Ok       Status = 0
	Warning  Status = 1
	Critical Status = 2
)

func Exit(status Status, msg string) {
	if msg != "" {
		msg = strings.TrimRight(msg, "\n")
		fmt.Printf("%s\n", msg)
	}
	os.Exit(int(status))
}
