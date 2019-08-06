package main

import (
	"testing"
)

func TestGetDSFromURLSigConfigFileName(t *testing.T) {
	expecteds := map[string]string{
		"url_sig_foo.config":                        "foo",
		"url_sig_.config":                           "",
		"url_sig.config":                            "",
		"url_sig_foo.conf":                          "",
		"url_sig_foo.confi":                         "",
		"url_sig_foo_bar_baz.config":                "foo_bar_baz",
		"url_sig_url_sig_foo_bar_baz.config.config": "url_sig_foo_bar_baz.config",
	}

	for fileName, expected := range expecteds {
		actual := GetDSFromURLSigConfigFileName(fileName)
		if expected != actual {
			t.Errorf("GetDSFromURLSigConfigFileName('%v') expected '%v' actual '%v'\n", fileName, expected, actual)
		}
	}
}
