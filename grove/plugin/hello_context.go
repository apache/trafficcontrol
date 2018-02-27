package plugin

import (
	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

func init() {
	AddPlugin(10000, Funcs{startup: helloCtxStart, afterRespond: helloCtxAfterResp})
}

func helloCtxStart(icfg interface{}, d StartupData) {
	*d.Context = 42
	log.Debugf("Hello World! Start set context: %+v\n", *d.Context)
}

func helloCtxAfterResp(icfg interface{}, d AfterRespondData) {
	ctx, ok := (*d.Context).(int)
	log.Debugf("Hello World! After Response got context: %+v %+v\n", ok, ctx)
}
