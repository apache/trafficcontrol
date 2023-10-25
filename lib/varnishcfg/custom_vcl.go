package varnishcfg

import (
	"strings"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
)

func (v VCLBuilder) configureCustomVCL(vclFile *vclFile) {
	params := t3cutil.FilterParams(v.toData.ServerParams, "default.vcl", "", "", "")
	for _, param := range params {
		if param.Name == "import" {
			vclFile.imports = append(vclFile.imports, param.Value)
			continue
		}
		// TODO: support loading vcl files too with `include`?? i.e. `include "custom.vcl";`
		lines := strings.Split(param.Value, "\n")
		vclFile.subroutines[param.Name] = append(vclFile.subroutines[param.Name], lines...)
	}
}
