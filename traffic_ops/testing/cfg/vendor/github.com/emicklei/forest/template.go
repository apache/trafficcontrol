package forest

import (
	"bytes"
	"text/template"
)

// ProcessTemplate creates a new text Template and executes it using the provided value.
// Returns the string result of applying this template.
// Failures in the template are reported using t.
func ProcessTemplate(t T, templateContent string, value interface{}) string {
	tmp, err := template.New("temporary").Parse(templateContent)
	if err != nil {
		logfatal(t, sfatalf("failed to parse:%v", err))
		return ""
	}
	var buf bytes.Buffer
	err = tmp.Execute(&buf, value)
	if err != nil {
		logfatal(t, sfatalf("failed to execute template:%v", err))
		return ""
	}
	return buf.String()
}
