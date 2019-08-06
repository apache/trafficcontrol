package main

import (
	"testing"
)

func TestGetDSFromURISigningConfigFileName(t *testing.T) {
	expecteds := map[string]string{
		"uri_signing_foo.config":                            "foo",
		"uri_signing_.config":                               "",
		"uri_signing.config":                                "",
		"uri_signing_foo.conf":                              "",
		"uri_signing_foo.confi":                             "",
		"uri_signing_foo_bar_baz.config":                    "foo_bar_baz",
		"uri_signing_uri_signing_foo_bar_baz.config.config": "uri_signing_foo_bar_baz.config",
	}

	for fileName, expected := range expecteds {
		actual := GetDSFromURISigningConfigFileName(fileName)
		if expected != actual {
			t.Errorf("GetDSFromURLSigConfigFileName('%v') expected '%v' actual '%v'\n", fileName, expected, actual)
		}
	}
}
