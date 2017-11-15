package forest

import (
	"fmt"
	"os"
	"testing"

	"github.com/wsxiaoys/terminal/color"
)

// TerminalColorsEnabled can be changed to disable the use of terminal coloring.
// One usecase is to add a command line flag to your test that controls its value.
//
//	func init() {
//		flag.BoolVar(&forest.TerminalColorsEnabled, "color", true, "want colors?")
//	}
//
//	go test -color=false
var TerminalColorsEnabled = true

// Check for presence of the TERMCOLORS environment variable to set the TerminalColorsEnabled setting.
func init() {
	TerminalColorsEnabled = os.Getenv("TERMCOLORS") != "false"
}

// ErrorColorSyntaxCode requires the syntax defined on https://github.com/wsxiaoys/terminal/blob/master/color/color.go .
// Set to an empty string to disable coloring.
var ErrorColorSyntaxCode = "@{wR}"

// FatalColorSyntaxCode requires the syntax defined on https://github.com/wsxiaoys/terminal/blob/master/color/color.go .
// Set to an empty string to disable coloring.
var FatalColorSyntaxCode = "@{wR}"

func serrorf(format string, args ...interface{}) string {
	return Scolorf(ErrorColorSyntaxCode, format, args...)
}

func sfatalf(format string, args ...interface{}) string {
	return Scolorf(FatalColorSyntaxCode, format, args...)
}

// Scolorf returns a string colorized for terminal output using the syntaxCode (unless that's empty).
// Requires the syntax defined on https://github.com/wsxiaoys/terminal/blob/master/color/color.go .
func Scolorf(syntaxCode string, format string, args ...interface{}) string {
	plainFormatted := fmt.Sprintf(format, args...)
	if len(syntaxCode) > 0 && TerminalColorsEnabled {
		// cannot pass the code as a string param
		return color.Sprintf(syntaxCode+"\n %s", plainFormatted)
	}
	return plainFormatted
}

// Errorf calls Error on t with a colorized message
func Errorf(t *testing.T, format string, args ...interface{}) {
	t.Error(Scolorf(ErrorColorSyntaxCode, format, args...))
}

// Fatalf calls Fatal on t with a colorized message
func Fatalf(t *testing.T, format string, args ...interface{}) {
	t.Fatal(Scolorf(FatalColorSyntaxCode, format, args...))
}
