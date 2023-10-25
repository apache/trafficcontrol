package varnishcfg

import (
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func TestConfigureCustomVCL(t *testing.T) {
	vb := NewVCLBuilder(&t3cutil.ConfigData{
		ServerParams: []tc.ParameterV50{
			{ConfigFile: "default.vcl", Name: "import", Value: "std"},
			{ConfigFile: "default.vcl", Name: "vcl_recv", Value: "set req.url = std.querysort(req.url);"},
			{
				ConfigFile: "default.vcl",
				Name:       "vcl_deliver",
				Value:      "if (req.status >= 400 && req.status <= 500) {\n\tset req.status = 404;\n}",
			},
		},
	})

	vclFile := newVCLFile(defaultVCLVersion)
	vb.configureCustomVCL(&vclFile)

	expectedVCLFile := newVCLFile(defaultVCLVersion)
	expectedVCLFile.imports = append(expectedVCLFile.imports, "std")
	expectedVCLFile.subroutines["vcl_recv"] = []string{
		"set req.url = std.querysort(req.url);",
	}
	expectedVCLFile.subroutines["vcl_deliver"] = []string{
		"if (req.status >= 400 && req.status <= 500) {",
		"	set req.status = 404;",
		"}",
	}

	if !reflect.DeepEqual(vclFile, expectedVCLFile) {
		t.Errorf("got %v want %v", vclFile, expectedVCLFile)
	}
}
