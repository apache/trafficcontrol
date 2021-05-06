package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall" // TODO change to x/unix ?
)

var commands = map[string]struct{}{
	"apply":  struct{}{},
	"update": struct{}{},
}

const ExitCodeNoCommand = 1
const ExitCodeUnknownCommand = 2
const ExitCodeCommandErr = 3
const ExitCodeExeErr = 4

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "no command\n") // TODO print usage
		os.Exit(ExitCodeNoCommand)
	}

	cmd := os.Args[1]
	if _, ok := commands[cmd]; !ok {
		fmt.Fprintf(os.Stderr, "unknown command\n") // TODO print usage
		os.Exit(ExitCodeUnknownCommand)
	}

	app := "t3c-" + cmd
	args := append([]string{app}, os.Args[2:]...)

	ex, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting application information: "+err.Error()+"\n")
		os.Exit(ExitCodeExeErr)
	}
	dir := filepath.Dir(ex)
	appDir := filepath.Join(dir, app) // TODO use path, not exact dir of this exe

	env := os.Environ()

	if err := syscall.Exec(appDir, args, env); err != nil {
		fmt.Fprintf(os.Stderr, "error executing sub-command: "+err.Error()+"\n")
		os.Exit(ExitCodeCommandErr)
	}
}
