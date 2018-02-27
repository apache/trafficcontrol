package plugin

import (
	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

func init() {
	AddPlugin(10000, Funcs{startup: hello})
}

func hello(icfg interface{}, d StartupData) {
	log.Errorf("Hello World! I'm a startup plugin! We're starting with %v bytes!\n", d.Config.CacheSizeBytes)
}
